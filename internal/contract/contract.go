// Package contract loads and validates the application's OpenAPI source document.
package contract

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
)

// Load rejects invalid documents and external references to keep validation local.
func Load(ctx context.Context, document []byte) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.Context = ctx
	spec, err := loader.LoadFromData(document)
	if err != nil {
		return nil, err
	}
	if err := spec.Validate(ctx); err != nil {
		return nil, err
	}
	return spec, nil
}
