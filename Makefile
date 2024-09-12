# Stop execution when a command fails
.ONESHELL:
.SHELLFLAGS = -ec

# Automatically determine the module path using go list
MODULE_PATH := $(shell go list -m -f {{.Path}})
COMMON_PACKAGE := $(MODULE_PATH)/common

# Define the output directory for binaries
BIN_DIR := bin

# Define the output binary names
DAEMON_BINARY := $(BIN_DIR)/pinner
CLI_BINARY := $(BIN_DIR)/pinner-cli
LC_BINARY := $(BIN_DIR)/light-client

# Define the source files
DAEMON_SOURCE := cmd/pinner/main.go
CLI_SOURCE := cmd/pinner-cli/main.go
LC_SOURCE := cmd/light-client/main.go

# Define the directories containing Go files
GO_DIRS := api cmd/pinner cmd/pinner-cli common internal

# Retrieve version and git commit hash
VERSION := $(shell git describe --tags --always --dirty)
GIT_COMMIT := $(shell git rev-parse HEAD)

# Default target to build all binaries
.PHONY: all
all: fmt vet staticcheck test build

# Create the bin directory if it doesn't exist
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Build both the daemon and CLI tool
.PHONY: build
build: $(BIN_DIR) buid-pinner build-cli build-light

# Build the daemon binary with a custom name
.PHONY: buid-pinner
buid-pinner: $(BIN_DIR)
	go build $(EXTRA_TAGS) -ldflags "-s -w -X $(COMMON_PACKAGE).BinaryName=pinner -X $(COMMON_PACKAGE).Version=$(VERSION) -X $(COMMON_PACKAGE).GitCommit=$(GIT_COMMIT)" -o $(DAEMON_BINARY) $(DAEMON_SOURCE)

# Build the CLI tool binary with a custom name
.PHONY: build-cli
build-cli: $(BIN_DIR)
	go build -ldflags "-s -w -X $(COMMON_PACKAGE).BinaryName=pinner-cli -X $(COMMON_PACKAGE).Version=$(VERSION) -X $(COMMON_PACKAGE).GitCommit=$(GIT_COMMIT)" -o $(CLI_BINARY) $(CLI_SOURCE)

# Build Light client
.PHONY: build-light
build-light: $(BIN_DIR)
	go build -ldflags "-s -w -X $(COMMON_PACKAGE).BinaryName=light-client -X $(COMMON_PACKAGE).Version=$(VERSION) -X $(COMMON_PACKAGE).GitCommit=$(GIT_COMMIT)" -o $(LC_BINARY) $(LC_SOURCE)

# Run tests
.PHONY: test
test:
	go test ./...

# Run tests with coverage report
.PHONY: test-cover
test-cover:
	go test -cover ./...

# Format the code
.PHONY: fmt
fmt:
	go fmt ./...

# Run static analysis (vet)
.PHONY: vet
vet:
	go vet ./...

# Tidy up module dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Clean up binaries and other generated files
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	rm -f coverage.out

# Generate code (if applicable)
.PHONY: generate
generate:
	go generate ./...

# Install dependencies
.PHONY: deps
deps:
	go mod download

.PHONY: staticcheck
staticcheck:
	staticcheck ./...