include .env
export

DB_URL := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

migrate-file:
	migrate create -ext sql -dir migrations/ -seq actor_table

migrate-force:
	migrate -path migrations -database "$(DB_URL)" -verbose force 1

migrate_up:
	migrate -path migrations/ -database "$(DB_URL)" up

migrate_version:
	migrate -path migrations/ -database "$(DB_URL)" version

migrate_down:
	migrate -path migrations/ -database "$(DB_URL)" down

run-local:
	go run cmd/movie-app/main.go