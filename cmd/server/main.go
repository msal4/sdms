package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/msal4/sdms"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Store struct {
	subjects []sdms.Subject
}

func (s *Store) GetSubjects() []sdms.Subject {
	return s.subjects
}

func (s *Store) AddSubject(subject sdms.Subject) {
}

const (
	pgUser     = "testuser"
	pgPassword = "testpassword"
	pgDb       = "testdb"
)

func main() {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image: "postgres",
		Env: map[string]string{
			"POSTGRES_USER":     pgUser,
			"POSTGRES_PASSWORD": pgPassword,
			"POSTGRES_DB":       pgDb,
		},
		WaitingFor:   wait.ForLog("Database connection is ready"),
		ExposedPorts: []string{"5432/tcp"},
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("couldnt start postgres container, %v", err)
	}

	p, _ := pgContainer.MappedPort(ctx, "5432")

	fmt.Printf("db is running on port %s\n", p.Port())

	defer pgContainer.Terminate(ctx)

	store := &Store{}

	srv := sdms.NewServer(store)

	log.Fatal(http.ListenAndServe(":5000", srv))
}
