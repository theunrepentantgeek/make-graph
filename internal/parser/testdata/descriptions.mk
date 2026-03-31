build: deps ## Build the binary
test-unit: ## Run unit tests
	go test ./...
clean: # Not a description
	rm -rf build/
deploy: deps ## Deploy to production # inline note
