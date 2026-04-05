.PHONY: help migrate-up migrate-down migrate-steps migrate-create migrate-force

# Colors for output
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

migrate-up: ## Apply all pending migrations
	@echo "$(GREEN)Applying all pending migrations...$(NC)"
	go run cmd/migrate/main.go -direction up

migrate-down: ## Rollback all migrations
	@echo "$(RED)Rolling back all migrations...$(NC)"
	go run cmd/migrate/main.go -direction down

migrate-steps: ## Apply N steps (usage: make migrate-steps steps=2)
	@echo "$(GREEN)Applying $(steps) migration steps...$(NC)"
	go run cmd/migrate/main.go -direction steps -steps $(steps)

migrate-create: ## Create new migration (usage: make migrate-create name=add_column_to_users)
	@echo "$(GREEN)Creating migration: $(name)$(NC)"
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	migrate create -ext sql -dir migrations -seq $(name)

migrate-force: ## Force set migration version (usage: make migrate-force version=1)
	@echo "$(RED)Forcing migration version to $(version)$(NC)"
	@echo "This can be dangerous! Make sure you know what you're doing."
	go run cmd/migrate/main.go -direction force -version $(version)

migrate-status: ## Show current migration status
	@echo "$(GREEN)Migration status:$(NC)"
	@psql $(DATABASE_URL) -c "SELECT version, dirty FROM schema_migrations;"

.PHONY: run-api
run-api: ## Run the API server
	@echo "$(GREEN)Starting API server...$(NC)"
	go run cmd/api/main.go