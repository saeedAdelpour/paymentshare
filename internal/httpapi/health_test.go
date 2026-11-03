package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saeedAdelpour/paymentshare/internal/httpapi"
)

func TestHealthDoesNotRequireDatabase(t *testing.T) {
	response := httptest.NewRecorder()
	httpapi.New(nil).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if response.Code != http.StatusOK || response.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected health response: %d %s", response.Code, response.Body.String())
	}
}

func TestReadinessWithoutDatabaseReturnsSafeUnavailableResponse(t *testing.T) {
	response := httptest.NewRecorder()
	httpapi.New(nil).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if response.Code != http.StatusServiceUnavailable || response.Body.String() != "{\"status\":\"unavailable\"}\n" {
		t.Fatalf("unexpected readiness response: %d %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type: %s", response.Header().Get("Content-Type"))
	}
}
