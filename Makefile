build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auction cmd/auction/main.go

run:
	docker compose up

down:
	docker compose down

test:
	go test -v internal/infra/database/auction/create_auction_test.go

.PHONY: build run down test
