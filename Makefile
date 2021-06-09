
migrate-up:
	./migrate.sh up

migrate-down:
	./migrate.sh down

test:
	go test -v ./...

start:
	./start.sh
