package modules

import (
	"net/http"
)

type BodyLimitConfig struct {
	MaxBytes int64
}

type BodyLimitHandler struct {
	Config BodyLimitConfig
	Next   http.Handler
}

func NewBodyLimitHandler(config BodyLimitConfig, next http.Handler) *BodyLimitHandler {
	return &BodyLimitHandler{
		Config: config,
		Next:   next,
	}
}

func (bl BodyLimitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		r.Body = http.MaxBytesReader(w, r.Body, bl.Config.MaxBytes)
	}

	bl.Next.ServeHTTP(w, r)
}
