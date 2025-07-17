# Build the application binary
build:
	go build -o bin/app ./cmd/app

# Run all tests with coverage
test:
	go test -v -race -coverprofile=coverage.out ./...

# Deploy to production environment
# Depends on build and test
deploy: build test
	./scripts/deploy.sh production

# f:name=program f:verb=run
# f:desc="Run main.go"
target:
	go run main.go

# Should be skipped
.PHONY: build test clean

# Should be skipped
%.o: %.c
	gcc -c $< -o $@