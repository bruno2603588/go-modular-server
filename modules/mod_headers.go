package modules

import (
	"net/http"
)

type HeadersConfig struct {
	CustomHeaders map[string]string
}

type HeadersHandler struct {
	Config HeadersConfig
	Next   http.Handler
}

func NewHeadersHandler(config HeadersConfig, next http.Handler) *HeadersHandler {
	return &HeadersHandler{
		Config: config,
		Next:   next,
	}
}

func (h HeadersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, value := range h.Config.CustomHeaders {
		w.Header().Set(key, value)
	}

	h.Next.ServeHTTP(w, r)
}
