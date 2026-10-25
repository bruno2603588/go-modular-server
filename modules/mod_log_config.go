package modules

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += int64(n)
	return n, err
}

type LogHandler struct {
	Next http.Handler
}

func NewLogHandler(next http.Handler) *LogHandler {
	return &LogHandler{Next: next}
}

func (l LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rw := &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}

	l.Next.ServeHTTP(rw, r)

	duration := time.Since(start)

	clientIP := r.RemoteAddr

	log.Printf("[ACCESS] IP: %-15s | Method: %-6s | Path: %-25s | Status: %3d | Size: %5d bytes | Duration: %v",
		clientIP,
		r.Method,
		r.URL.Path,
		rw.status,
		rw.size,
		duration,
	)
}
