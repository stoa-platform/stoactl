.PHONY: build build-agent build-all install clean test lint fmt vet run

# Variables
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
LDFLAGS_CLI=-ldflags "-s -w -X github.com/stoa-platform/stoa-go/internal/cli/cmd.Version=$(VERSION) -X github.com/stoa-platform/stoa-go/internal/cli/cmd.Commit=$(COMMIT)"
LDFLAGS_AGENT=-ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT)"

# Build stoactl
build:
	go build $(LDFLAGS_CLI) -o bin/stoactl ./cmd/stoactl

# Build stoa-agent
build-agent:
	go build $(LDFLAGS_AGENT) -o bin/stoa-agent ./cmd/stoa-agent

# Install to GOPATH/bin
install:
	go install $(LDFLAGS_CLI) ./cmd/stoactl
	go install $(LDFLAGS_AGENT) ./cmd/stoa-agent

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
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS_CLI) -o bin/stoactl-darwin-amd64 ./cmd/stoactl
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS_CLI) -o bin/stoactl-darwin-arm64 ./cmd/stoactl
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS_CLI) -o bin/stoactl-linux-amd64 ./cmd/stoactl
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS_CLI) -o bin/stoactl-linux-arm64 ./cmd/stoactl
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS_CLI) -o bin/stoactl-windows-amd64.exe ./cmd/stoactl
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS_AGENT) -o bin/stoa-agent-darwin-amd64 ./cmd/stoa-agent
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS_AGENT) -o bin/stoa-agent-darwin-arm64 ./cmd/stoa-agent
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS_AGENT) -o bin/stoa-agent-linux-amd64 ./cmd/stoa-agent
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS_AGENT) -o bin/stoa-agent-linux-arm64 ./cmd/stoa-agent

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
