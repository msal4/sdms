#!/bin/sh

. .env

migrate -path migrations -database "$DB_URL" $1
