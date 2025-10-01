.PHONY: build run clean test install deps

# Build variables
BINARY_NAME=multichat
MAIN_PATH=.
BUILD_DIR=.

# Build the application
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build with optimizations for production
build-prod:
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application (WhatsApp with debug logging)
run:
	go run $(MAIN_PATH) --messenger whatsapp --device device.db --log-level debug

# Run with custom parameters
run-custom:
	go run $(MAIN_PATH) $(ARGS)

# Clean build artifacts
clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f *.db *.db-shm *.db-wal

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install the binary to GOPATH/bin
install:
	go install

# Cross-compile for different platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-prod    - Build with optimizations"
	@echo "  run           - Run with default WhatsApp settings"
	@echo "  run-custom    - Run with custom args (use ARGS='...')"
	@echo "  clean         - Remove build artifacts and databases"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  build-all     - Cross-compile for all platforms"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"

