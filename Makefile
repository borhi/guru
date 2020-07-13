MIGRATE=./migrate.darwin-amd64 -path db/migrations -database mongodb://mongo:mongo@localhost:27017/guru
TEST_MIGRATE=./migrate.darwin-amd64 -path db/test_migrations -database mongodb://mongo:mongo@localhost:27017/test_guru

.PHONY: build start stop down migrate-up migrate-down test

build: ## Build docker containers
	docker-compose build

start: ## Start docker containers
	docker-compose up -d

stop: ## Stop docker containers
	docker-compose stop

down: ## Down docker containers
	docker-compose down --volumes

migrate-up: ## Run migrations
	$(MIGRATE) up

migrate-down: ## Rollback migrations
	$(MIGRATE) down

test: ## Request test
	docker-compose -f docker-compose.test.yml up -d --build
	sleep 2
	$(TEST_MIGRATE) up
	go test -v ./...
	docker-compose -f docker-compose.test.yml down --volumes