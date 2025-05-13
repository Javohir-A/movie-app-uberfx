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
IMAGE_NAME = javohirgo/movie_app
TAG = v1.0.0

build-image:
	docker build -t $(IMAGE_NAME):$(TAG) .

push-image:
	docker push $(IMAGE_NAME):$(TAG)

run-docker:
	docker compose up -d


swag-v1: ### swag init
	swag init -g internal/router/router.go

run:
	docker network inspect movies-network >/dev/null 2>&1 || docker network create movies-network
	docker compose up --build
