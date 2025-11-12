.PHONY: help
help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Install all required tools for development
	@echo "Installing development tools..."
	@go mod download
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo "Setup complete!"

.PHONY: tools
tools: ## Verify all required tools are installed
	@echo "Checking required tools..."
	@command -v goose >/dev/null 2>&1 || { echo "goose not found. Run 'make setup'"; exit 1; }
	@command -v gofumpt >/dev/null 2>&1 || { echo "gofumpt not found. Run 'make setup'"; exit 1; }
	@command -v oapi-codegen >/dev/null 2>&1 || { echo "oapi-codegen not found. Run 'make setup'"; exit 1; }
	@echo "All tools are installed!"

.PHONY: tidy
tidy: ## Tidy and verify go.mod
	@go mod tidy
	@go mod verify

.PHONY: fmt
fmt: ## Format code with gofumpt
	@gofumpt -l -w .

.PHONY: fmt-check
fmt-check: ## Check if code is formatted correctly
	@test -z "$$(gofumpt -l .)" || (echo "Code is not formatted. Run 'make fmt'"; gofumpt -l .; exit 1)

.PHONY: vet
vet: ## Run go vet for static analysis
	@go vet ./...

.PHONY: lint
lint: fmt-check vet ## Run all linting checks (format check + vet)

.PHONY: test
test: ## Run all tests
	@go test ./... -v

.PHONY: test-cover
test-cover: ## Run tests with coverage
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: build
build: ## Build the application binary
	@echo "Building athylps..."
	@go build -o bin/athylps cmd/athylps/main.go
	@echo "Binary created at bin/athylps"

.PHONY: run
run: ## Run the application
	@go run cmd/athylps/main.go

.PHONY: clean
clean: ## Clean build artifacts and coverage reports
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Cleaned build artifacts"

.PHONY: migrations-up
migrations-up: ## Run database migrations up
	@go run cmd/migrate/main.go up

.PHONY: migrations-down
migrations-down: ## Rollback last database migration
	@go run cmd/migrate/main.go down

.PHONY: migrations-status
migrations-status: ## Check database migration status
	@go run cmd/migrate/main.go status

.PHONY: migration
migration: ## Create a new migration file (usage: make migration name=create_users_table)
	@test -n "$(name)" || (echo "name is required. Usage: make migration name=create_users_table"; exit 1)
	@goose -dir ./migrations create $(name) sql

.PHONY: generate-api-models
generate-api-models: # (Re)Generate api request/response models based on api specification
	@oapi-codegen --config=api/oapi-codegen.yaml api/openapi.yaml

.PHONY: check
check: lint test ## Run all checks (lint + test)

.PHONY: ci
ci: tools tidy lint test build ## Run all CI checks (formatting, linting, tests, build)

.DEFAULT_GOAL := help
