APP_NAME=dbz
MAIN=./cmd/main.go
OUTPUT=bin/${APP_NAME}

IMAGE=${APP_NAME}:latest

.PHONY: all dev build test tidy lint clean

dev:
	@echo "Starting development mode..."
	@go tool air

build:
	@echo "Building binary..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(OUTPUT) $(MAIN)
	@echo "Binary generated at: $(OUTPUT)"

test:
	@echo "Running tests with coverage..."
	@go test ./... -cover -race

lint:
	@echo "Running golangci-lint..."

tidy:
	@echo "Tidying modules..."
	@go mod tidy

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin
