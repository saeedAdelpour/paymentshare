// Package testutil contains shared HTTP test harness behavior.
package testutil

import (
	"bytes"
	"fmt"
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
	if err := CheckExchange(request, response); err != nil {
		t.Fatal(err)
	}
}

// CheckExchange reports contract mismatches, including undeclared redirect statuses.
func CheckExchange(request *http.Request, response *http.Response) error {
	spec, err := contract.Load(request.Context(), api.Document)
	if err != nil {
		return err
	}
	router, err := legacy.NewRouter(spec)
	if err != nil {
		return err
	}
	route, params, err := router.FindRoute(request)
	if err != nil {
		return err
	}
	if route.Operation.Responses.Status(response.StatusCode) == nil && route.Operation.Responses.Default() == nil {
		return fmt.Errorf("response status %d is not declared for %s %s", response.StatusCode, request.Method, request.URL.Path)
	}
	input := &openapi3filter.RequestValidationInput{Request: request, PathParams: params, Route: route}
	if err := openapi3filter.ValidateRequest(request.Context(), input); err != nil {
		return err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err := response.Body.Close(); err != nil {
		return err
	}
	response.Body = io.NopCloser(bytes.NewReader(body))
	if err := openapi3filter.ValidateResponse(request.Context(), &openapi3filter.ResponseValidationInput{
		RequestValidationInput: input,
		Status:                 response.StatusCode,
		Header:                 response.Header,
		Body:                   io.NopCloser(bytes.NewReader(body)),
		Options:                &openapi3filter.Options{IncludeResponseStatus: true},
	}); err != nil {
		return err
	}
	return nil
}
