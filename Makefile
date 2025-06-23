include tilt/Makefile

.PHONY: check fmt vet test lint test-e2e

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
	golangci-lint --build-tags e2e run

# Run tests
test:
	go test ./...

test-e2e:
	go test -v -p 1 -count=1 -tags=e2e -timeout=10m ./e2e
