DB_URL=postgres://postgres:postgres@localhost:5432/books?sslmode=disable
MIGRATIONS_DIR=./db/migrations

migrate-up:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) up

migrate-down:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) down

sqlc-gen:
	sqlc generate

psql:
	psql $(DB_URL)

drop-db:
	dropdb --if-exists -U postgres books

up:
	docker-compose up -d

down:
	docker-compose down

seed:
	psql $(DB_URL) -f db/seed.sql

.PHONY: migrate-up migrate-down migrate-new migrate-force sqlc-gen create-db drop-db psql up down logs