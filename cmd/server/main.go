package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/msal4/sdms"
)

func main() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("problem establishing a db connection: %v", err)
	}

	store := sdms.NewPostgresStore(db)
	server := sdms.NewServer(store)

	log.Fatal(http.ListenAndServe(":5000", server))
}
