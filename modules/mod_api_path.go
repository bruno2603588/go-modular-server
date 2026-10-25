package modules

import (
	"net/http"
)

type APIConfig struct {
	EnableCORS     bool
	Prefix         string
	AllowedOrigins []string
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type APIHandler struct {
	Config APIConfig
	Router *http.ServeMux
}

func NewAPIHandler(config APIConfig) *APIHandler {
	return &APIHandler{
		Config: config,
		Router: http.NewServeMux(),
	}
}

func (h APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.Config.EnableCORS {
		origin := r.Header.Get("Origin")

		if len(h.Config.AllowedOrigins) > 0 {
			if h.isOriginAllowed(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	h.Router.ServeHTTP(w, r)
}

func (h APIHandler) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range h.Config.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
