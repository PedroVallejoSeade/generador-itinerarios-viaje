.PHONY: build test fmt vet check run

# Build the CLI binary
build:
	go build -o bin/citysearch ./cmd/citysearch

# Run the full test suite
test:
	go test ./...

# Format all Go sources
fmt:
	gofmt -w .

# Static analysis
vet:
	go vet ./...

# Format check + vet + tests
check:
	gofmt -l .
	go vet ./...
	go test ./...

# Build and run (usage: make run ARGS="paris")
run: build
	./bin/citysearch $(ARGS)
