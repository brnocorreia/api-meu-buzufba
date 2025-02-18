run:
	go run cmd/api/main.go

build:
	go build -o api cmd/api/main.go

mongo:
	docker compose up -d

mongo-stop:
	docker compose down

