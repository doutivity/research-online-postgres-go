POSTGRES_URI=postgresql://yaroslav:AnySecretPassword!!@localhost:5432/yaaws?sslmode=disable

env-up:
	docker-compose up -d

env-down:
	docker-compose down --remove-orphans -v

go-test:
	docker exec research-online-postgres-go-app go test ./... -v -count=1

docker-go-version:
	docker exec research-online-postgres-go-app go version
	docker exec research-online-postgres-go-app go run main.go

docker-pg-version:
	docker exec research-online-postgres-1 psql -U yaroslav -d yaaws -c "SELECT VERSION();"

go-bench:
	mkdir -p ./output/

	docker exec research-online-postgres-go-app go test ./... -v -run=$$^ -bench='TxLoopUpsert'    -benchmem -benchtime=1000x -count=10 | tee ./output/bench-go-1000x-tx-loop-upsert.txt
	docker exec research-online-postgres-go-app go test ./... -v -run=$$^ -bench='TxLoopUpdate'    -benchmem -benchtime=1000x -count=10 | tee ./output/bench-go-1000x-tx-loop-update.txt
	docker exec research-online-postgres-go-app go test ./... -v -run=$$^ -bench='BatchExecUpsert' -benchmem -benchtime=1000x -count=10 | tee ./output/bench-go-1000x-batch-exec-upsert.txt
	docker exec research-online-postgres-go-app go test ./... -v -run=$$^ -bench='BatchExecUpdate' -benchmem -benchtime=1000x -count=10 | tee ./output/bench-go-1000x-batch-exec-update.txt
	docker exec research-online-postgres-go-app go test ./... -v -run=$$^ -bench='UnnestUpsert'    -benchmem -benchtime=1000x -count=10 | tee ./output/bench-go-1000x-unnest-upsert.txt
	docker exec research-online-postgres-go-app go test ./... -v -run=$$^ -bench='UnnestUpdate'    -benchmem -benchtime=1000x -count=10 | tee ./output/bench-go-1000x-unnest-update.txt

	# go install golang.org/x/perf/cmd/benchstat@latest

	# benchstat ./output/bench-go-1000x-tx-loop-upsert.txt
	# benchstat ./output/bench-go-1000x-tx-loop-update.txt
	# benchstat ./output/bench-go-1000x-batch-exec-upsert.txt
	# benchstat ./output/bench-go-1000x-batch-exec-update.txt
	# benchstat ./output/bench-go-1000x-unnest-upsert.txt
	# benchstat ./output/bench-go-1000x-unnest-update.txt

	# benchstat ./output/bench-all.txt

test:
	make env-up
	make docker-go-version
	make docker-pg-version
	make migrate-up
	make go-test
	make env-down

# go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# sqlc generate
#
# alternative
# docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate
generate-sqlc:
	sqlc generate

# Creates new migration file with the current timestamp
# Example: make create-new-migration-file NAME=<name>
create-new-migration-file:
	$(eval NAME ?= noname)
	mkdir -p ./migrations/
	goose -dir ./migrations/ create $(NAME) sql

migrate-up:
	goose -dir ./migrations/ -table schema_migrations postgres $(POSTGRES_URI) up
migrate-redo:
	goose -dir ./migrations/ -table schema_migrations postgres $(POSTGRES_URI) redo
migrate-down:
	goose -dir ./migrations/ -table schema_migrations postgres $(POSTGRES_URI) down
migrate-reset:
	goose -dir ./migrations/ -table schema_migrations postgres $(POSTGRES_URI) reset
migrate-status:
	goose -dir ./migrations/ -table schema_migrations postgres $(POSTGRES_URI) status

install-sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest
