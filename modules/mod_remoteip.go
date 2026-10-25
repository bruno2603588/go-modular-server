package modules

import (
	"net"
	"net/http"
	"strings"
)

type RemoteIPConfig struct {
	HeaderField    string
	TrustedProxies []string
}

type RemoteIPHandler struct {
	Config RemoteIPConfig
	Next   http.Handler
}

func NewRemoteIPHandler(config RemoteIPConfig, next http.Handler) *RemoteIPHandler {
	if config.HeaderField == "" {
		config.HeaderField = "X-Forwarded-For"
	}
	return &RemoteIPHandler{
		Config: config,
		Next:   next,
	}
}

func (m *RemoteIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r, m.Config)

	if clientIP != "" {
		r.RemoteAddr = clientIP
	}

	m.Next.ServeHTTP(w, r)
}

func getClientIP(r *http.Request, config RemoteIPConfig) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	if len(config.TrustedProxies) > 0 && !isTrustedProxy(host, config.TrustedProxies) {
		return host
	}

	headerValue := r.Header.Get(config.HeaderField)
	if headerValue != "" {
		ips := strings.Split(headerValue, ",")
		realIP := strings.TrimSpace(ips[0])
		if net.ParseIP(realIP) != nil {
			return realIP
		}
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		if net.ParseIP(realIP) != nil {
			return realIP
		}
	}

	return host
}

func isTrustedProxy(ipStr string, trustedList []string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, trusted := range trustedList {
		if trusted == ipStr {
			return true
		}
		_, cidr, err := net.ParseCIDR(trusted)
		if err == nil && cidr.Contains(ip) {
			return true
		}
	}
	return false
}
