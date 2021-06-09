package sdms_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/msal4/sdms"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	pgImage    = "postgres:12.4-alpine"
	pgUser     = "postgres"
	pgPassword = "postgres"
	pgDB       = "testdb"

	port   = "5432/tcp"
	driver = "postgres"
)

func openTestDB() (*sql.DB, func(), error) {
	ctx := context.Background()

	env := map[string]string{
		"POSTGRES_USER":     pgUser,
		"POSTGRES_PASSWORD": pgPassword,
		"POSTGRES_DB":       pgDB,
	}

	dbURL := func(p nat.Port) string {
		return fmt.Sprintf("%s://%s:%s@localhost:%s/%s?sslmode=disable", driver, pgUser, pgPassword, p.Port(), pgDB)
	}

	req := testcontainers.ContainerRequest{
		Image:        pgImage,
		Env:          env,
		WaitingFor:   wait.ForSQL(nat.Port(port), driver, dbURL).Timeout(time.Second * 5),
		ExposedPorts: []string{port},
	}
	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("couldnt start postgres container: %v", err)
	}

	terminate := func() { pgC.Terminate(ctx) }

	p, err := pgC.MappedPort(ctx, nat.Port(port))
	if err != nil {
		terminate()
		return nil, nil, fmt.Errorf("couldnt get mapped port: %v", err)
	}

	db, err := sql.Open(driver, dbURL(p))
	if err != nil {
		terminate()
		return nil, nil, fmt.Errorf("problem connecting to postgres db: %v", err)
	}

	if err = db.Ping(); err != nil {
		terminate()
		return nil, nil, fmt.Errorf("ping to postgres db failed: %v", err)
	}

	return db, terminate, nil
}

func getMigrate(db *sql.DB) (*migrate.Migrate, error) {
	driverWithDBInstance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create postgres driver instance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", driver, driverWithDBInstance)
	if err != nil {
		return nil, fmt.Errorf("problem creating migrate instance: %v", err)
	}

	return m, nil
}

func migrateUp(db *sql.DB) error {
	m, err := getMigrate(db)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("problem migrating up: %v", err)
	}

	return nil
}

func migrateDown(db *sql.DB) error {
	m, err := getMigrate(db)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("problem migrating down: %v", err)
	}

	return nil
}

func TestPostgresStore_Subjects(t *testing.T) {
	db, terminate, err := openTestDB()
	if err != nil {
		t.Fatalf("failed to openen test db: %v", err)
	}
	t.Cleanup(terminate)

	store := sdms.NewPostgresStore(db)

	initDB := func(t *testing.T) {
		if err := migrateDown(db); err != nil {
			t.Fatalf("failed to migrate down: %v", err)
		}
		if err := migrateUp(db); err != nil {
			t.Fatalf("failed to run migrations: %v", err)
		}
	}

	t.Run("get subjects", func(t *testing.T) {
		initDB(t)

		subjects, err := store.GetSubjects()
		assertNoError(t, err)
		initialSubjectsLength := len(subjects)

		_, err = db.Exec(`
		insert into lecturers (name, image, password, about) values ($1, $2, $3, $4);
		`,
			"test subject",
			"test/image",
			"mypassword",
			"about me",
		)

		if err != nil {
			t.Fatalf("problem executing lecturer insert query: %v", err)
		}

		_, err = db.Exec(`
		insert into subjects (name, details, stage, lecturer_id, semester, syllabus) values ($1, $2, $3, $4, $5, $6);
		`,
			"test subject",
			"test details",
			1,
			1,
			sdms.FirstSemester,
			"test syllabus",
		)
		if err != nil {
			t.Fatalf("problem executing insert query: %v", err)
		}

		subjects, err = store.GetSubjects()

		assertNoError(t, err)

		assertSubjectsLength(t, subjects, initialSubjectsLength+1)
	})

	t.Run("add subject", func(t *testing.T) {
		initDB(t)

		subjects, err := store.GetSubjects()
		assertNoError(t, err)
		initialSubjectsLength := len(subjects)

		subject := sdms.Subject{
			Name:     "Geography",
			Details:  "Meta",
			Semester: sdms.FirstSemester,
			Syllabus: "Hi",
			Lecturer: &sdms.Lecturer{
				ID:    1874,
				Name:  "Sub",
				Image: "/images/2.png",
				About: "I'm a professor",
			},
		}

		err = store.AddSubject(&subject)
		assertNoError(t, err)

		subjects, err = store.GetSubjects()
		assertNoError(t, err)

		assertSubjectsLength(t, subjects, initialSubjectsLength+1)

		assertHasSubject(t, subjects, subject)
	})

	t.Run("remove subject", func(t *testing.T) {
		initDB(t)

		subjects, err := store.GetSubjects()
		assertNoError(t, err)
		initialSubjectsLength := len(subjects)

		subject := sdms.Subject{
			Name:     "Geography",
			Details:  "Meta",
			Semester: sdms.FirstSemester,
			Syllabus: "Hi",
			Lecturer: &sdms.Lecturer{
				ID:    1874,
				Name:  "Sub",
				Image: "/images/2.png",
				About: "I'm a professor",
			},
		}

		err = store.AddSubject(&subject)
		assertNoError(t, err)

		err = store.RemoveSubject(subject.ID)
		assertNoError(t, err)

		subjects, err = store.GetSubjects()
		assertNoError(t, err)

		assertSubjectsLength(t, subjects, initialSubjectsLength)

		assertNotHasSubject(t, subjects, subject)
	})

	t.Run("update subject", func(t *testing.T) {
		initDB(t)

		want := sdms.Subject{
			Name:     "Geography",
			Details:  "Meta",
			Semester: sdms.FirstSemester,
			Syllabus: "Hi",
			Lecturer: &sdms.Lecturer{
				ID:    1874,
				Name:  "Sub",
				Image: "/images/2.png",
				About: "I'm a professor",
			},
		}

		err = store.AddSubject(&want)
		assertNoError(t, err)

		want.Name = "New name"

		err = store.UpdateSubject(want)
		assertNoError(t, err)

		got, err := store.GetSubjectByID(want.ID)
		assertNoError(t, err)

		if got.Name != want.Name {
			t.Errorf("got subject name %q want %q", got.Name, want.Name)
		}

	})
}

func assertNotHasSubject(t testing.TB, subjects []sdms.Subject, want sdms.Subject) {
	t.Helper()
	var found *sdms.Subject
	for _, s := range subjects {
		if s.Name == want.Name && s.Details == want.Details && s.Semester == want.Semester && s.Syllabus == want.Syllabus {
			found = &s
		}
	}

	if found != nil {
		t.Fatalf("expected no subject got %#v", found)
	}

}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("did not expect an error: %v", err)
	}
}

func assertHasError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error")
	}
}
