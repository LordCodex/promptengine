BINARY_NAME=promptengine
VERSION?=$(shell cat VERSION 2>/dev/null || echo "0.1.0-alpha")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-X github.com/LordCodex/promptengine/internal/version.Version=$(VERSION) \
                  -X github.com/LordCodex/promptengine/internal/version.Commit=$(COMMIT) \
                  -X github.com/LordCodex/promptengine/internal/version.BuildDate=$(BUILD_DATE) -s -w"

.PHONY: all build clean test lint fmt run install release snapshot bench

all: fmt lint test build

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/promptengine/main.go

clean:
	rm -rf bin/ dist/ coverage.txt

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

bench:
	go test -run=^$$ -bench=. -benchmem ./tests/bench/...

lint:
	go vet ./...
	@if command -v staticcheck >/dev/null 2>&1; then staticcheck ./...; else echo "staticcheck not installed, skipping"; fi

fmt:
	go fmt ./...

run:
	go run cmd/promptengine/main.go $(ARGS)

install:
	go install $(LDFLAGS) ./cmd/promptengine

snapshot:
	@if command -v goreleaser >/dev/null 2>&1; then goreleaser release --snapshot --clean; else echo "goreleaser not installed, fallback build"; go build $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/promptengine/main.go; fi

release:
	@if command -v goreleaser >/dev/null 2>&1; then goreleaser release --clean; else echo "goreleaser not installed"; exit 1; fi
