# SDMS (Scientific Department Management System)
Proof of concept.

### Setup

- Install ![golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

- Create a `.env` file using `cp .env.example .env` and set `DB_URL`

- Run `make migrate-up` to run the migrations

- Run `make test` to run the tests

- Run `make start` to start the development server


