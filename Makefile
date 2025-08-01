DB_URL=postgres://postgres:postgres@localhost:5432/books?sslmode=disable
MIGRATIONS_DIR=./db/migrations

.PHONY: help migrate-up migrate-down migrate-force sqlc-gen psql up down build clean

help: ## Prints all targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

migrate-up: ## Migrate up database
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) up

migrate-down: ## Migrate down database
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) down

migrate-force: ## Force migrate database to version 1
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) force 1

sqlc-gen: ## Generate Go code from SQL using sqlc
	sqlc generate

psql: ## Open interactive psql session
	psql $(DB_URL)

up: ## Start services via docker-compose
	docker-compose up -d

down: ## Stop services via docker-compose
	docker-compose down

build: ## Compile the Go application
	go build -o bin/semantic-search ./cmd/main.go

clean: ## Clean build files (asks for confirmation)
	@read -p "Are you sure you want to delete ./bin? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "Cleaning..."; \
		rm -rf bin/; \
	else \
		echo "Aborted."; \
	fi