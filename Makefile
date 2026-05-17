include .env
export

.PHONY: run
run:
	go run ./cmd/sso/main.go -config_path=$(CONFIG_PATH)

.PHONY: migrate-up
migrate-up:
	go run ./cmd/migrator/main.go -command="up" -dir=$(MIGRATION_PATH) -config_path=$(CONFIG_PATH)

.PHONY: migrate-down
migrate-down:
	go run ./cmd/migrator/main.go -command="down" -dir=$(MIGRATION_PATH) -config_path=$(CONFIG_PATH)

.PHONY: migrate-refresh
migrate-refresh:
	go run ./cmd/migrator/main.go -command="refresh" -dir=$(MIGRATION_PATH) -config_path=$(CONFIG_PATH)

.PHONY: docker-build
docker-build:
	docker-compose build

.PHONY: docker-up
docker-up:
	docker-compose --env-file .env.docker up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: docker-down-volumes
docker-down-volumes:
	docker-compose down -v

.PHONY: docker-logs
docker-logs:
	docker-compose logs -f

.PHONY: docker-status
docker-status:
	docker-compose ps

.PHONY: docker-restart
docker-restart:
	docker-compose restart app

.PHONY: docker-clean
docker-clean:
	docker-compose down -v