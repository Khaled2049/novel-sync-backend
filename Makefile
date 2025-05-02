# Makefile

# Ensure .env file exists and load variables
# This guards against errors if .env is missing and makes variables available
# Using include and export allows variables to be directly used by commands
# If .env doesn't exist, it includes an empty file, preventing errors.
# Use := for immediate evaluation, avoids potential issues with complex commands.
-include .env
export $(shell sed 's/=.*//' .env)

# --- Variables ---
# Default migration name if MIGRATION_NAME is not provided via command line
MIGRATION_NAME ?= new_migration
# Database URL from environment (loaded from .env)
DB_URL := ${DATABASE_URL}
# Migrations directory
MIGRATIONS_DIR := migrations

# Check if migrate CLI is installed
MIGRATE_BIN := $(shell command -v migrate 2> /dev/null)

# --- Helper Targets ---
.PHONY: help check_migrate check_db_url

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

check_migrate: ## Check if migrate CLI is installed
ifndef MIGRATE_BIN
	$(error "migrate CLI not found in PATH. Install it: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest")
endif

check_db_url: ## Check if DATABASE_URL is set
ifndef DB_URL
	$(error "DATABASE_URL is not set. Please ensure it is defined in your .env file.")
endif

# --- Migration Targets ---
.PHONY: migrate-create migrate-up migrate-down migrate-down-all migrate-force migrate-status db-start db-stop db-logs

db-start: ## Start the postgres docker container
	@echo "Starting PostgreSQL container..."
	docker-compose up -d db

db-stop: ## Stop the postgres docker container
	@echo "Stopping PostgreSQL container..."
	docker-compose down

db-logs: ## View logs for the postgres container
	@echo "Tailing PostgreSQL container logs..."
	docker-compose logs -f db

# Example: make migrate-create name=add_users_table
migrate-create: check_migrate ## Create new up/down migration files (e.g., make migrate-create name=add_index)
	@read -p "Enter migration name (e.g., add_users_index): " name; \
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

migrate-up: check_migrate check_db_url ## Apply all pending UP migrations
	@echo "Running UP migrations..."
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) up

migrate-down: check_migrate check_db_url ## Revert the last applied migration
	@echo "Running DOWN migration (reverting last)..."
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) down 1

migrate-down-all: check_migrate check_db_url ## Revert all applied migrations (Use with caution!)
	@echo "Running DOWN migration (reverting ALL)..."
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) down -all

# Use force carefully, only if the migration state is known to be incorrect
migrate-force: check_migrate check_db_url ## Force migration version (e.g., make migrate-force version=1)
	@read -p "Enter migration version number to force: " version; \
	echo "Forcing migration version to $$version (use with caution)..."; \
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) force $$version

migrate-status: check_migrate check_db_url ## Show current migration status
	@echo "Checking migration status..."
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) version