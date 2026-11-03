.PHONY: dev deps tools run migrate db-up db-down generate check-generated check-contract test-contract-tools check-assets fmt check-format lint vuln build check test test-integration ci

GO_DIRS = api cmd internal integration migrations

deps:
	go mod download
	go -C tools mod download

tools: bin/oapi-codegen bin/gofumpt bin/golangci-lint bin/govulncheck

bin/oapi-codegen: tools/go.mod tools/go.sum
	go -C tools build -mod=readonly -o $(CURDIR)/$@ github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen

bin/gofumpt: tools/go.mod tools/go.sum
	go -C tools build -mod=readonly -o $(CURDIR)/$@ mvdan.cc/gofumpt

bin/golangci-lint: tools/go.mod tools/go.sum
	go -C tools build -mod=readonly -o $(CURDIR)/$@ github.com/golangci/golangci-lint/v2/cmd/golangci-lint

bin/govulncheck: tools/go.mod tools/go.sum
	go -C tools build -mod=readonly -o $(CURDIR)/$@ golang.org/x/vuln/cmd/govulncheck

dev:
	@sh scripts/with-env.sh sh -c 'export SWAGGER_ENABLED="$${SWAGGER_ENABLED:-true}"; exec go run ./cmd/api'

run:
	go run ./cmd/api

migrate:
	@sh scripts/with-env.sh go run ./cmd/migrate $(if $(TO),-to $(TO),)

db-up:
	docker compose up -d --wait postgres

db-down:
	docker compose down

generate: bin/oapi-codegen
	bin/oapi-codegen -config api/oapi-codegen.yaml -o internal/httpapi/generated.go api/openapi.yaml

check-generated: bin/oapi-codegen
	@tmp=$$(mktemp); trap 'rm -f "$$tmp"' EXIT; \
	bin/oapi-codegen -config api/oapi-codegen.yaml -o "$$tmp" api/openapi.yaml && \
	diff -u internal/httpapi/generated.go "$$tmp"

check-contract:
	go run ./cmd/contract $(if $(SPEC),-file $(SPEC),)

fmt: bin/gofumpt
	bin/gofumpt -w $(GO_DIRS)

check-format: bin/gofumpt
	@files=$$(bin/gofumpt -l $(GO_DIRS)) || exit $$?; \
	if [ -n "$$files" ]; then printf 'Run make fmt:\n%s\n' "$$files"; exit 1; fi

lint: bin/golangci-lint
	bin/golangci-lint run --timeout=5m

vuln: bin/govulncheck
	bin/govulncheck -tags=integration ./...

build:
	go build ./...

test-contract-tools: bin/oapi-codegen
	sh scripts/test-contract-checks.sh

check-assets:
	cd internal/httpapi/docs && sha256sum -c SHA256SUMS

check: check-assets check-format check-contract check-generated test-contract-tools build lint vuln

test:
	go test -race ./...

test-integration:
	go test -race -tags=integration ./integration -count=1

# CI and local full verification use the same command. A real PostgreSQL URL is required.
ci: check
	go test -race -tags=integration ./... -count=1
