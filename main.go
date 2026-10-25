package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"http-server/modules"
)

func main() {
	handler := setupRouter()

	server := &http.Server{
		Addr:              ":8000",
		Handler:           handler,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       40 * time.Second,
		WriteTimeout:      20 * time.Second,
	}

	log.Println("Servidor en marcha en http://localhost:8000")
	log.Fatal(server.ListenAndServe())
}

func setupRouter() http.Handler {
	mux := http.NewServeMux()

	apiModule := modules.NewAPIHandler(modules.APIConfig{
		EnableCORS:     true,
		Prefix:         "/api",
		AllowedOrigins: loadCorsOrigins(),
	})

	apiModule.Router.HandleFunc("/api/public/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Servidor operando con Rate Limiting activo"}`))
	})

	secureMux := http.NewServeMux()
	secureMux.HandleFunc("/api/private/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Panel privado"}`))
	})

	authzModule := modules.NewAuthzHandler(modules.AuthzConfig{
		ValidTokens: loadAllowedTokens(),
	}, secureMux)

	apiModule.Router.Handle("/api/private/", authzModule)
	mux.Handle("/api/", apiModule)

	dirModule := modules.NewDirHandler(modules.DirConfig{
		StaticDir:      "./dist",
		DirectoryIndex: "index.html",
	})
	mux.Handle("/", dirModule)

	logModule := modules.NewLogHandler(mux)

	headersModule := modules.NewHeadersHandler(modules.HeadersConfig{
		CustomHeaders: map[string]string{
			"Strict-Transport-Security": "max-age=63072000; includeSubDomains; preload",
			"X-Content-Type-Options":    "nosniff",
			"X-Frame-Options":           "DENY",
			"Content-Security-Policy":   "default-src 'self'; script-src 'self';",
		},
	}, logModule)

	bodyLimitModule := modules.NewBodyLimitHandler(modules.BodyLimitConfig{
		MaxBytes: 1024 * 1024, // 1 MB
	}, headersModule)

	rateLimitModule := modules.NewRateLimitHandler(modules.RateLimitConfig{
		MaxRequests: 5,
		Window:      10 * time.Second,
	}, bodyLimitModule)

	remoteIPModule := modules.NewRemoteIPHandler(modules.RemoteIPConfig{
		HeaderField:    "X-Forwarded-For",
		TrustedProxies: []string{"127.0.0.1", "::1"},
	}, rateLimitModule)

	return remoteIPModule
}

func loadAllowedTokens() []string {
	raw := os.Getenv("API_TOKENS")
	if raw == "" {
		return nil
	}
	tokens := strings.Split(raw, ",")
	for i, t := range tokens {
		tokens[i] = strings.TrimSpace(t)
	}
	return tokens
}

func loadCorsOrigins() []string {
	raw := os.Getenv("CORS_ORIGINS")
	if raw == "" {
		return []string{"http://localhost:5173", "http://localhost:3000"}
	}
	if raw == "*" {
		return []string{"*"}
	}
	origins := strings.Split(raw, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}
	return origins
}
