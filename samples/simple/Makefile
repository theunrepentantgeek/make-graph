.PHONY: all build test clean deploy lint test-integration

all: build test ## Build and test everything

build: deps ## Build the binary
	go build -o bin/app ./cmd/app

deps: ## Install dependencies
	go mod download

test: build ## Run tests
	go test ./...

test-integration: test ## Run integration tests
	go test -tags=integration ./...

clean: ## Clean build artifacts
	rm -rf bin/

deploy: build test ## Deploy to production
	@$(MAKE) -C deploy all

lint: ## Run linter
	golangci-lint run
