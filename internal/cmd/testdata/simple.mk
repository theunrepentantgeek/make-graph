.PHONY: all build test clean

all: build test ## Build and test everything

build: ## Build the binary
	go build ./...

test: build ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -rf build/
