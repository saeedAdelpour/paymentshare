// Package testutil contains shared HTTP test harness behavior.
package testutil

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/saeedAdelpour/paymentshare/api"
	"github.com/saeedAdelpour/paymentshare/internal/contract"
)

// ValidateExchange checks both sides of an HTTP exchange against the source
// contract and restores the response body for the caller's behavioral assertions.
func ValidateExchange(t testing.TB, request *http.Request, response *http.Response) {
	t.Helper()
	spec, err := contract.Load(context.Background(), api.Document)
	if err != nil {
		t.Fatal(err)
	}
	router, err := legacy.NewRouter(spec)
	if err != nil {
		t.Fatal(err)
	}
	route, params, err := router.FindRoute(request)
	if err != nil {
		t.Fatal(err)
	}
	input := &openapi3filter.RequestValidationInput{Request: request, PathParams: params, Route: route}
	if err := openapi3filter.ValidateRequest(request.Context(), input); err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatal(err)
	}
	response.Body = io.NopCloser(bytes.NewReader(body))
	if err := openapi3filter.ValidateResponse(request.Context(), &openapi3filter.ResponseValidationInput{
		RequestValidationInput: input,
		Status:                 response.StatusCode,
		Header:                 response.Header,
		Body:                   io.NopCloser(bytes.NewReader(body)),
	}); err != nil {
		t.Fatal(err)
	}
}
