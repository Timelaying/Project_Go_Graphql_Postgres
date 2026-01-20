run:
	go run ./cmd/api

gen:
	go run github.com/99designs/gqlgen generate

up:
	docker compose up -d

down:
	docker compose down -v
