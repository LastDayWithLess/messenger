
BINARY_NAME=messenger
DB_URL=postgres://user123:qwerty123@localhost:5432/messenger?sslmode=require
MIGRATION_DIR=migrations

.PHONY: help run build migration migrate-up migrate-status migrate-down test-unit

help:
	@echo "Доступные команды:"
	@echo "  make build       - Собрать приложение"
	@echo "  make run         - Запустить приложение"
	@echo "Миграции:"
	@echo "  make migration   - Создать новую миграцию"
	@echo "  make migrate-up  - Применить миграции"
	@echo "  make migrate-down - Откатить последнюю миграцию"
	@echo "  make migrate-status - Показать статус миграций"

run:
	go run cmd/app/main.go

build:
	go build -o bin/$(BINARY_NAME) cmd/app/main.go

migration:
	@read -p "Migration name:" name; \
	goose -dir $(MIGRATION_DIR) create $$name sql

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

default: help


test-unit:
	go test ./internal/transport -v