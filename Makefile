

# Директория, где хранятся файлы миграций
MIGRATIONS_DIR := ./schemas

# Создание новой миграции
.PHONY: migrate-create
migrate-create:
	@goose -dir=$(MIGRATIONS_DIR) create $(name) sql

# Применение всех миграций
.PHONY: migrate-up
migrate-up:
	@goose -dir=$(MIGRATIONS_DIR) postgres $(DB_CONN_STRING) up

# Откат последней миграции
.PHONY: migrate-down
migrate-down:
	goose -dir=$(MIGRATIONS_DIR) postgres $(DB_CONN_STRING) down