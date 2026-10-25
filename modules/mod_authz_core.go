package modules

import (
	"encoding/json"
	"net/http"
	"strings"
)

type AuthzConfig struct {
	ValidTokens []string
	RequireFunc func(r *http.Request) bool
}

type AuthzHandler struct {
	Config AuthzConfig
	Next   http.Handler
}

func NewAuthzHandler(config AuthzConfig, next http.Handler) *AuthzHandler {
	return &AuthzHandler{
		Config: config,
		Next:   next,
	}
}

func (a AuthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.Config.RequireFunc != nil {
		if !a.Config.RequireFunc(r) {
			sendAuthzError(w, http.StatusForbidden, "Acceso prohibido: Reglas de autorización no cumplidas")
			return
		}
	}

	if len(a.Config.ValidTokens) > 0 {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			sendAuthzError(w, http.StatusUnauthorized, "No autorizado: Falta el token de acceso")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		tokenAuthorized := false

		for _, vt := range a.Config.ValidTokens {
			if vt == token {
				tokenAuthorized = true
				break
			}
		}

		if !tokenAuthorized {
			sendAuthzError(w, http.StatusForbidden, "Acceso denegado: Token inválido o expirado")
			return
		}
	}

	a.Next.ServeHTTP(w, r)
}

func sendAuthzError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
