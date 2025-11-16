BIN=server

include .env
export $(shell sed 's/=.*//' .env)

TEST_COMPOSE=docker-compose.test.yml

.PHONY: gen
gen:
	go tool oapi-codegen -config oapi.cfg.yml openapi.yml

.PHONY: run
run: stop
	@echo "Starting server"
	docker-compose up --build -d
	@echo "Server is up on http://localhost:${SERVER_PORT}"

.PHONY: stop
stop:
	@echo "Stopping server"
	docker-compose down

.PHONY: run-local
run-local:
	@echo "Running server locally"
	go run ./cmd/$(BIN)

.PHONY: test-e2e
test-e2e:
	@echo "Starting E2E tests"
	@trap 'echo docker-compose -f $(TEST_COMPOSE) down -v' EXIT; \
		docker-compose -f $(TEST_COMPOSE) up -d postgres-test; \
		sleep 3; \
		echo "Running migrations on test DB..."; \
		goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:${DB_PORT}/testdb?sslmode=disable" up; \
		echo "Running tests..."; \
		go test ./tests/e2e -count=1 -v; \
		test_exit_code=$$?; \
		echo "Resetting test database..."; \
		goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:${DB_PORT}/testdb?sslmode=disable" reset; \
		docker-compose -f $(TEST_COMPOSE) down postgres-test;
		exit $$test_exit_code

.PHONY: logs
logs:
	docker-compose logs -f

.PHONY: clean
clean: stop
	docker-compose -f $(TEST_COMPOSE) down -v 2>/dev/null || true
	docker system prune -f

.PHONY: migrate-up
migrate-up:
	GOOSE_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
	goose -dir db/migrations postgres up

.PHONY: migrate-down
migrate-down:
	GOOSE_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
	goose -dir db/migrations postgres down
