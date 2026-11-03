// Package server composes the HTTP and persistence adapters.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/saeedAdelpour/paymentshare/internal/httpapi"
	"github.com/saeedAdelpour/paymentshare/internal/postgres"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func LoadConfig() (Config, error) {
	cfg := Config{HTTPAddr: os.Getenv("HTTP_ADDR"), DatabaseURL: os.Getenv("DATABASE_URL")}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = "127.0.0.1:8082"
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if _, _, err := net.SplitHostPort(cfg.HTTPAddr); err != nil {
		return Config{}, errors.New("HTTP_ADDR must have host:port format")
	}
	return cfg, nil
}

// Run serves until cancellation, drains in-flight requests, and closes the pool.
func Run(ctx context.Context, cfg Config, log *slog.Logger) error {
	startup, cancel := context.WithTimeout(ctx, 5*time.Second)
	_, pool, err := postgres.Open(startup, cfg.DatabaseURL)
	cancel()
	if err != nil {
		// Driver errors may contain credentials. Never log or return the raw error here.
		return errors.New("PostgreSQL startup check failed; check DATABASE_URL, credentials, TLS, and database availability")
	}
	defer func() {
		if err := pool.Close(); err != nil {
			log.Error("database pool close failed")
		}
	}()
	if err := ctx.Err(); err != nil {
		return err
	}
	listener, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		return fmt.Errorf("listen HTTP: %w", err)
	}
	srv := &http.Server{
		Handler:           httpapi.New(pool),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}
	stopped := make(chan error, 1)
	go func() { stopped <- srv.Serve(listener) }()
	log.Info("HTTP server listening", "address", listener.Addr().String())
	select {
	case err := <-stopped:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			_ = srv.Close()
			return fmt.Errorf("shutdown HTTP: %w", err)
		}
		if err := <-stopped; !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		log.Info("HTTP server stopped")
		return nil
	}
}
