//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/saeedAdelpour/paymentshare/internal/httpapi"
	"github.com/saeedAdelpour/paymentshare/internal/postgres"
	"github.com/saeedAdelpour/paymentshare/internal/server"
	"github.com/saeedAdelpour/paymentshare/internal/testutil"
)

func TestReadinessReflectsRealPostgresConnectivity(t *testing.T) {
	dsn, _ := newDatabase(t)
	runMigration(t, dsn, "1", true)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, pool, err := postgres.Open(ctx, dsn)
	if err != nil {
		t.Fatal("could not connect to test PostgreSQL (details suppressed to protect credentials)")
	}
	t.Cleanup(func() { _ = pool.Close() })
	handler := httpapi.New(pool, false)
	request := func(path string, code int, body string) {
		t.Helper()
		response := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		handler.ServeHTTP(response, req)
		result := response.Result()
		testutil.ValidateExchange(t, req, result)
		_ = result.Body.Close()
		if response.Code != code || response.Body.String() != body {
			t.Fatalf("%s: got %d %s", path, response.Code, response.Body.String())
		}
	}
	request("/readyz", http.StatusOK, "{\"status\":\"ok\"}\n")
	if err := pool.Close(); err != nil {
		t.Fatal("closing test pool failed")
	}
	request("/readyz", http.StatusServiceUnavailable, "{\"status\":\"unavailable\"}\n")
	request("/healthz", http.StatusOK, "{\"status\":\"ok\"}\n")
}

// readyLog observes the bound address without choosing a potentially occupied port.
type readyLog struct{ address chan string }

func (l readyLog) Write(data []byte) (int, error) {
	var entry struct {
		Address string `json:"address"`
	}
	if err := json.Unmarshal(data, &entry); err != nil {
		return 0, err
	}
	if entry.Address != "" {
		l.address <- entry.Address
	}
	return len(data), nil
}

func TestServerStartsAndShutsDownWithPostgres(t *testing.T) {
	dsn, _ := newDatabase(t)
	runMigration(t, dsn, "1", true)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := readyLog{address: make(chan string, 1)}
	done := make(chan error, 1)
	go func() {
		done <- server.Run(ctx, server.Config{HTTPAddr: "127.0.0.1:0", DatabaseURL: dsn}, slog.New(slog.NewJSONHandler(log, nil)))
	}()
	var address string
	select {
	case address = <-log.address:
	case err := <-done:
		t.Fatalf("server did not start: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("server startup timed out")
	}
	client := &http.Client{Timeout: 3 * time.Second}
	response, err := client.Get("http://" + address + "/readyz")
	if err != nil {
		t.Fatal(err)
	}
	testutil.ValidateExchange(t, response.Request, response)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("readiness status: %d", response.StatusCode)
	}
	runMigration(t, dsn, "0", true)
	unavailable, err := client.Get("http://" + address + "/readyz")
	if err != nil {
		t.Fatal(err)
	}
	testutil.ValidateExchange(t, unavailable.Request, unavailable)
	_ = unavailable.Body.Close()
	if unavailable.StatusCode != http.StatusServiceUnavailable {
		t.Fatal("running server did not detect schema reversal")
	}
	runMigration(t, dsn, "1", true)
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(12 * time.Second):
		t.Fatal("server shutdown timed out")
	}
	response, err = client.Get("http://" + address + "/healthz")
	if err == nil {
		_ = response.Body.Close()
		t.Fatal("server still accepts requests after shutdown")
	}
}
