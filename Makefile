ENTRYPOINT := main.go
OUTPUT := bc

# Default target to build the application
all: build

# Build the Go application
build:
	go build -o $(OUTPUT) $(ENTRYPOINT)


count: build
	./$(OUTPUT) count

sync: build
	./$(OUTPUT) sync

# Clean up build artifacts
clean:
	rm -f $(OUTPUT)

# Run tests
test:
	go test ./...

# Phony targets (targets that are not actual files)
.PHONY: all build run clean test