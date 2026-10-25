package modules

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Handler dummy de prueba para el final de la cadena de middlewares
var dummyNext = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
})

// ==========================================
// 1. PRUEBAS PARA MOD_REMOTEIP
// ==========================================
func TestRemoteIPHandler(t *testing.T) {
	// Configuramos el módulo con un proxy confiable: 127.0.0.1
	config := RemoteIPConfig{
		HeaderField:    "X-Forwarded-For",
		TrustedProxies: []string{"127.0.0.1"},
	}

	handler := NewRemoteIPHandler(config, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Comprobamos que RemoteAddr haya sido sobrescrito con la IP del header
		if r.RemoteAddr != "203.0.113.195" {
			t.Errorf("Esperaba la IP 203.0.113.195, pero se obtuvo %s", r.RemoteAddr)
		}
		w.WriteHeader(http.StatusOK)
	}))

	// Simulamos una petición desde el proxy confiable 127.0.0.1 con la cabecera X-Forwarded-For
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.195, 10.0.0.1")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Código de estado incorrecto: obtuve %d, esperaba %d", rec.Code, http.StatusOK)
	}
}

// ==========================================
// 2. PRUEBAS PARA MOD_AUTHZ_CORE
// ==========================================
func TestAuthzHandler(t *testing.T) {
	config := AuthzConfig{
		ValidTokens: []string{"token-valido-abc"},
	}

	handler := NewAuthzHandler(config, dummyNext)

	// Caso A: Petición sin token (Debe retornar 401 Unauthorized)
	t.Run("Sin token de autorización", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/private", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Esperaba 401 Unauthorized, obtuve %d", rec.Code)
		}
	})

	// Caso B: Petición con token inválido (Debe retornar 403 Forbidden)
	t.Run("Token inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/private", nil)
		req.Header.Set("Authorization", "Bearer token-falso-xyz")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Esperaba 403 Forbidden, obtuve %d", rec.Code)
		}
	})

	// Caso C: Petición con token válido (Debe retornar 200 OK)
	t.Run("Token válido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/private", nil)
		req.Header.Set("Authorization", "Bearer token-valido-abc")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Esperaba 200 OK, obtuve %d", rec.Code)
		}
	})
}

// ==========================================
// 3. PRUEBAS PARA MOD_RATELIMIT
// ==========================================
func TestRateLimitHandler(t *testing.T) {
	// Límite estricto: máximo 2 peticiones por cada 5 segundos
	config := RateLimitConfig{
		MaxRequests: 2,
		Window:      5 * time.Second,
	}

	handler := NewRateLimitHandler(config, dummyNext)
	clientIP := "192.168.1.50:5000"

	// Petición 1 (Debe pasar)
	req1 := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	req1.RemoteAddr = clientIP
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Errorf("Petición 1 falló: obtuve %d", rec1.Code)
	}

	// Petición 2 (Debe pasar)
	req2 := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	req2.RemoteAddr = clientIP
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Errorf("Petición 2 falló: obtuve %d", rec2.Code)
	}

	// Petición 3 (Debe ser bloqueada con 429 Too Many Requests)
	req3 := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	req3.RemoteAddr = clientIP
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusTooManyRequests {
		t.Errorf("Esperaba código 429 Too Many Requests por exceder el límite, obtuve %d", rec3.Code)
	}
}
