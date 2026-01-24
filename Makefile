.PHONY: build install clean test lint fmt vet run

# Variables
BINARY_NAME=stoactl
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
LDFLAGS=-ldflags "-s -w -X github.com/stoa-platform/stoactl/internal/cmd.Version=$(VERSION) -X github.com/stoa-platform/stoactl/internal/cmd.Commit=$(COMMIT)"

# Build
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/stoactl

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/stoactl

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf dist/

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run the CLI
run:
	go run ./cmd/stoactl $(ARGS)

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/stoactl
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/stoactl
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/stoactl
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 ./cmd/stoactl
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/stoactl

# Release with goreleaser (dry-run)
release-dry:
	goreleaser release --snapshot --clean

# Check goreleaser config
release-check:
	goreleaser check

# Dependencies
deps:
	go mod download
	go mod tidy
