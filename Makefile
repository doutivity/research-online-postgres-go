POSTGRES_URI=postgresql://yaroslav:AnySecretPassword!!@localhost:5432/yaaws?sslmode=disable

env-up:
	docker-compose up -d

test:
	docker exec research-online-postgres-go-app go test ./... -v -count=1

generate-sqlc:
	# go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	sqlc generate
	# alternative
	# docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

migrate-create:
	mkdir -p ./storage/schema
	goose -dir ./storage/schema -table schema_migrations postgres $(POSTGRES_URI) create $(MIGRATION_NAME) sql

migrate-up:
	goose -dir ./storage/schema -table schema_migrations postgres $(POSTGRES_URI) up
migrate-redo:
	goose -dir ./storage/schema -table schema_migrations postgres $(POSTGRES_URI) redo
migrate-down:
	goose -dir ./storage/schema -table schema_migrations postgres $(POSTGRES_URI) down
migrate-reset:
	goose -dir ./storage/schema -table schema_migrations postgres $(POSTGRES_URI) reset
migrate-status:
	goose -dir ./storage/schema -table schema_migrations postgres $(POSTGRES_URI) status

env-down:
	docker-compose down --remove-orphans -v
