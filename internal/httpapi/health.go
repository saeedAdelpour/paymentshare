package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type healthServer struct{ pool *sql.DB }

// New serves the health contract. A nil pool reports unavailable readiness.
func New(pool *sql.DB) http.Handler {
	router := chi.NewRouter()
	return HandlerFromMux(&healthServer{pool: pool}, router)
}

func (s *healthServer) GetHealth(w http.ResponseWriter, r *http.Request) {
	writeHealth(w, http.StatusOK, Ok)
}

func (s *healthServer) GetReadiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if s.pool == nil || s.pool.PingContext(ctx) != nil {
		writeHealth(w, http.StatusServiceUnavailable, Unavailable)
		return
	}
	writeHealth(w, http.StatusOK, Ok)
}

func writeHealth(w http.ResponseWriter, code int, status HealthStatus) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	// A failed response write means the client has disconnected; no retry is safe.
	_ = json.NewEncoder(w).Encode(Health{Status: status})
}
