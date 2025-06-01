.PHONY: build run clean test screenshots

# Detect OS and set binary extension
ifeq ($(OS),Windows_NT)
    BINARY_EXT=.exe
else
    BINARY_EXT=
endif

# Binary name
BINARY_NAME=timetracker$(BINARY_EXT)

# Build the application
build:
	go mod download
	go build -o $(BINARY_NAME) ./cmd/timetracker

# Run the application
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f current_session.json
	rm -rf screenshots/

# Run tests
test:
	go test -v ./...

# Generate screenshots (useful for documentation)
screenshots:
	go test -v ./test -run TestGenerateScreenshots

# Install dependencies
deps:
	go mod download
	go mod tidy

# Default target
all: build 