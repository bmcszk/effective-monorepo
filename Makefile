include tilt/Makefile

.PHONY: check fmt vet test lint

# Default target
check: fmt vet lint test

# Format code
fmt:
	go fmt ./...

# Vet code for issues
vet:
	go vet ./...

# Run golangci-lint
lint:
	golangci-lint run

# Run tests
test:
	go test ./...
