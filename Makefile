.PHONY: help up down build logs health clean ps restart test-flow

COMPOSE := docker compose

##@ General
help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1m%-20s %s\033[0m\n\n", "Target", "Description"} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

##@ Development
up: ## Build and start all services in the background
	$(COMPOSE) up -d --build

down: ## Stop and remove all containers
	$(COMPOSE) down

build: ## Build all service images without starting
	$(COMPOSE) build

restart: ## Restart all running services
	$(COMPOSE) restart

ps: ## Show status of all services
	$(COMPOSE) ps

##@ Logs
logs: ## Tail logs from ALL services
	$(COMPOSE) logs -f

logs-user: ## Tail user-service logs
	$(COMPOSE) logs -f user-service

logs-post: ## Tail post-service logs
	$(COMPOSE) logs -f post-service

logs-notif: ## Tail notification-service logs
	$(COMPOSE) logs -f notification-service

##@ Health & Testing
health: ## Check health endpoints of all services
	@echo "\n\033[1m─── User Service (8081) ───\033[0m"
	@curl -sf http://localhost:8081/api/v1/health 2>/dev/null | python3 -m json.tool || echo "  ❌ DOWN"
	@echo "\n\033[1m─── Post Service (8082) ───\033[0m"
	@curl -sf http://localhost:8082/api/v1/health 2>/dev/null | python3 -m json.tool || echo "  ❌ DOWN"
	@echo "\n\033[1m─── Notification Service (8083) ───\033[0m"
	@curl -sf http://localhost:8083/api/v1/health 2>/dev/null | python3 -m json.tool || echo "  ❌ DOWN"

test-flow: ## Run a quick end-to-end smoke test
	@echo "\n\033[1m=== 1. Creating a user ===\033[0m"
	@curl -s -X POST http://localhost:8081/api/v1/users \
		-H "Content-Type: application/json" \
		-d '{"username":"kamal","email":"kamal@pulse.dev","bio":"DevOps Engineer"}' | python3 -m json.tool
	@sleep 1
	@echo "\n\033[1m=== 2. Listing all users ===\033[0m"
	@curl -s http://localhost:8081/api/v1/users | python3 -m json.tool
	@echo "\n\033[1m=== 3. Creating a post (triggers notification) ===\033[0m"
	@USER_ID=$$(curl -s http://localhost:8081/api/v1/users | python3 -c "import sys,json; print(json.load(sys.stdin)['data'][0]['id'])") && \
	curl -s -X POST http://localhost:8082/api/v1/posts \
		-H "Content-Type: application/json" \
		-d "{\"user_id\":\"$$USER_ID\",\"content\":\"First post on Pulse! 🚀\"}" | python3 -m json.tool
	@sleep 1
	@echo "\n\033[1m=== 4. Checking notifications ===\033[0m"
	@curl -s http://localhost:8083/api/v1/notifications | python3 -m json.tool
	@echo "\n\033[32m✅ End-to-end flow complete! Check 'make logs-notif' for event delivery.\033[0m"

##@ Cleanup
clean: ## Remove all containers, volumes, and locally built images
	$(COMPOSE) down -v --rmi local
	@echo "\033[32m✅ Cleaned up all Pulse resources.\033[0m"
