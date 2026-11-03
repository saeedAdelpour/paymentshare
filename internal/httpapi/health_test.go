package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saeedAdelpour/paymentshare/internal/testutil"

	"github.com/saeedAdelpour/paymentshare/internal/httpapi"
)

func TestHealthDoesNotRequireDatabase(t *testing.T) {
	response := httptest.NewRecorder()
	httpapi.New(nil, false).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	validateContract(t, "/healthz", response)
	if response.Code != http.StatusOK || response.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected health response: %d %s", response.Code, response.Body.String())
	}
}

func TestReadinessWithoutDatabaseReturnsSafeUnavailableResponse(t *testing.T) {
	response := httptest.NewRecorder()
	httpapi.New(nil, false).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	validateContract(t, "/readyz", response)
	if response.Code != http.StatusServiceUnavailable || response.Body.String() != "{\"status\":\"unavailable\"}\n" {
		t.Fatalf("unexpected readiness response: %d %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type: %s", response.Header().Get("Content-Type"))
	}
}

func TestDocumentationIsExplicitlyEnabled(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		handler := httpapi.New(nil, enabled)
		for _, path := range []string{"/docs/", "/docs/swagger-ui.css", "/docs/swagger-ui-bundle.js", "/openapi.yaml"} {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
			want := http.StatusNotFound
			if enabled {
				want = http.StatusOK
			}
			if response.Code != want {
				t.Fatalf("docs enabled=%v %s: got %d, want %d", enabled, path, response.Code, want)
			}
		}
	}
}

func validateContract(t *testing.T, path string, response *httptest.ResponseRecorder) {
	t.Helper()
	result := response.Result()
	defer func() { _ = result.Body.Close() }()
	testutil.ValidateExchange(t, httptest.NewRequest(http.MethodGet, path, nil), result)
}
