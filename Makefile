.PHONY: deps run generate check-generated check test test-integration

OAPI_CODEGEN = go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.1

deps:
	go mod download

run:
	go run ./cmd/api

generate:
	$(OAPI_CODEGEN) -config api/oapi-codegen.yaml -o internal/httpapi/generated.go api/openapi.yaml

check-generated:
	@tmp=$$(mktemp); trap 'rm -f "$$tmp"' EXIT; \
	$(OAPI_CODEGEN) -config api/oapi-codegen.yaml -o "$$tmp" api/openapi.yaml && \
	diff -u internal/httpapi/generated.go "$$tmp"

check: check-generated
	@test -z "$$(gofmt -l cmd internal integration)" || { gofmt -l cmd internal integration; exit 1; }
	go build ./...
	go vet ./...

test:
	go test -race ./...

test-integration:
	go test -race -tags=integration ./integration -count=1
