package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/saeedAdelpour/paymentshare/migrations"
)

type healthServer struct{ pool *sql.DB }

// New serves the health contract. A nil pool reports unavailable readiness.
func New(pool *sql.DB, docsEnabled bool) http.Handler {
	router := chi.NewRouter()
	if docsEnabled {
		registerDocumentation(router)
	}
	return HandlerFromMux(&healthServer{pool: pool}, router)
}

func (s *healthServer) GetHealth(w http.ResponseWriter, r *http.Request) {
	writeHealth(w, http.StatusOK, Health{Status: Ok})
}

func (s *healthServer) GetReadiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if s.pool == nil || migrations.CheckSchema(ctx, s.pool) != nil {
		writeHealth(w, http.StatusServiceUnavailable, UnavailableHealth{Status: Unavailable})
		return
	}
	writeHealth(w, http.StatusOK, Health{Status: Ok})
}

func writeHealth(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	// A failed response write means the client has disconnected; no retry is safe.
	_ = json.NewEncoder(w).Encode(body)
}
