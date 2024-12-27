DB_URL=postgres://root:root@localhost:5433/tengahin?sslmode=disable

.PHONY: run
run:
	docker compose up -d
	air

.PHONY: migrate
migrate:
	migrate create -ext sql -dir db/migration -seq $(word 2, $(MAKECMDGOALS))

.PHONY: migrateup
migrateup:
	migrate -path db/migration -database $(DB_URL) -verbose up

.PHONY: migratedown
migratedown:
	migrate -path db/migration -database $(DB_URL) -verbose down

.PHONY: sqlc
sqlc:
	sqlc generate