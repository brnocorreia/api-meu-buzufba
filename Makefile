run:
	go run cmd/api/main.go

build:
	go build -o api cmd/api/main.go

docker-build:
	docker compose build

docker:
	docker compose up -d

docker-stop:
	docker compose down

