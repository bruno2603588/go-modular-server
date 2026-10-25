package modules

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimitConfig struct {
	MaxRequests int
	Window      time.Duration
}

type clientStats struct {
	count     int
	resetTime time.Time
}

type RateLimitHandler struct {
	Config   RateLimitConfig
	Next     http.Handler
	mu       sync.Mutex
	visitors map[string]*clientStats
}

func NewRateLimitHandler(config RateLimitConfig, next http.Handler) *RateLimitHandler {
	limiter := &RateLimitHandler{
		Config:   config,
		Next:     next,
		visitors: make(map[string]*clientStats),
	}

	go limiter.cleanupVisitors()

	return limiter
}

func (rl *RateLimitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := extractIP(r.RemoteAddr)

	rl.mu.Lock()
	stats, exists := rl.visitors[ip]
	now := time.Now()

	if !exists || now.After(stats.resetTime) {
		rl.visitors[ip] = &clientStats{
			count:     1,
			resetTime: now.Add(rl.Config.Window),
		}
		rl.mu.Unlock()
	} else {
		stats.count++
		if stats.count > rl.Config.MaxRequests {
			rl.mu.Unlock()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests) // Código HTTP 429
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Demasiadas solicitudes. Límite de tasa excedido (Rate Limit Exceeded).",
			})
			return
		}
		rl.mu.Unlock()
	}

	rl.Next.ServeHTTP(w, r)
}

func extractIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

func (rl *RateLimitHandler) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		now := time.Now()
		for ip, stats := range rl.visitors {
			if now.After(stats.resetTime.Add(time.Minute)) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
