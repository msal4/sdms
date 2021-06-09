
migrate-up:
	./migrate.sh up

migrate-down:
	./migrate.sh down

test:
	go test -v ./...

start:
	go run cmd/server/main.go
