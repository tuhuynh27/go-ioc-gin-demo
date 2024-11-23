.PHONY: run test clean build

# Default binary output
BINARY_NAME=app

# Build the application
build:
	go build -o ${BINARY_NAME} main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -f coverage.out
	rm -f coverage.html

# Install dependencies
deps:
	go mod download

# Run linter
lint:
	go vet ./...
	golangci-lint run

# Migrate database
migrate:
	go run migrations/runner.go
