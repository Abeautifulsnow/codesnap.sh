build:
	go build -o ./bin/server .

run: build
	 ./bin/server

up:
	docker-compose up -d && docker compose logs -f app

test:
	go test -v ./...

fmt:
	go fmt ./...