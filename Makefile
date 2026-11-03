.PHONY: dev deps run migrate check test

dev:
	docker compose up -d

deps:
	go mod download

run:
	go run ./cmd/api

migrate:
	go run ./cmd/migrate

check:
	go build ./...
	go vet ./...

test:
	go test -race ./...
