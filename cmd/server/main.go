package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/msal4/sdms"
)

func loadEnv() error {
	if err := godotenv.Load("../../.env"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			return fmt.Errorf("Error loading .env file: %v", err)
		}
	}

	return nil
}

func main() {
	if err := loadEnv(); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("problem establishing a db connection: %v", err)
	}

	store := sdms.NewPostgresStore(db)
	server := sdms.NewServer(store)

	log.Println("Running on http://0.0.0.0:5000")

	log.Fatal(http.ListenAndServe(":5000", server))
}
