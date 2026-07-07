package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type response map[string]any

var (
	appName = requireEnv("APP_NAME")
	apiKey  = requireEnv("API_KEY")
	goEnv   = requireEnv("GO_ENV")
	release = requireEnv("RELEASE")
	port    = requireEnv("PORT")
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /secure", requireAPIKey(secureHandler))
	mux.HandleFunc("POST /echo", requireAPIKey(echoHandler))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("%s listening on port %s", appName, port)
	log.Fatal(server.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, response{
		"app":     appName,
		"release": release,
		"message": "Go API deployed through Portiq auto-detection (pushed)",
		"env": response{
			"apiKeyConfigured": apiKey != "",
			"goEnv":            goEnv,
			"port":             port,
		},
		"routes": []string{
			"GET /",
			"GET /health",
			"GET /secure",
			"POST /echo",
		},
	})
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, response{
		"ok":        true,
		"service":   appName,
		"release":   release,
		"checkedAt": time.Now().UTC().Format(time.RFC3339),
	})
}

func secureHandler(w http.ResponseWriter, _ *http.Request) {
	preview := ""
	if len(apiKey) >= 4 {
		preview = apiKey[:4] + "..."
	}

	writeJSON(w, http.StatusOK, response{
		"message":       "Go protected route reached",
		"secretPreview": preview,
	})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		payload = map[string]any{}
	}

	writeJSON(w, http.StatusOK, response{
		"app":      appName,
		"received": payload,
		"release":  release,
	})
}

func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if apiKey == "" {
			writeJSON(w, http.StatusInternalServerError, response{
				"error": "API_KEY is not configured",
			})
			return
		}

		if r.Header.Get("x-api-key") != apiKey {
			writeJSON(w, http.StatusUnauthorized, response{
				"error": "Invalid or missing API key",
			})
			return
		}

		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}

	return value
}
