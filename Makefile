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
