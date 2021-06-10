package sdms

import (
	"database/sql"
	"fmt"
)

var ErrNotFound = fmt.Errorf("object not found")

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) GetSubjects() ([]Subject, error) {
	rows, err := s.db.Query(`SELECT id, name, details, semester, stage, syllabus, lecturer_id FROM subjects;`)
	if err != nil {
		return nil, fmt.Errorf("error querying subjects: %v", err)
	}

	subjects := []Subject{}
	for rows.Next() {
		subject := Subject{
			Lecturer: &Lecturer{},
		}

		err := rows.Scan(&subject.ID, &subject.Name, &subject.Details, &subject.Semester, &subject.Stage, &subject.Syllabus, &subject.Lecturer.ID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning row: %v", err)
		}

		subjects = append(subjects, subject)
	}
	return subjects, nil
}

func (s *PostgresStore) GetSubjectByID(id int) (*Subject, error) {
	row := s.db.QueryRow("SELECT id, name, details, semester, stage, syllabus, lecturer_id FROM subjects where id = $1", id)

	subject := Subject{
		Lecturer: &Lecturer{},
	}

	err := row.Scan(&subject.ID, &subject.Name, &subject.Details, &subject.Semester, &subject.Stage, &subject.Syllabus, &subject.Lecturer.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("encountered error while scanning row: %v", err)
	}

	return &subject, nil
}

func (s *PostgresStore) AddSubject(subject *Subject) error {
	if subject.Lecturer == nil {
		subject.Lecturer = &Lecturer{}
	}
	row := s.db.QueryRow(`INSERT INTO subjects (name, details, semester, stage, syllabus, lecturer_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`,
		subject.Name, subject.Details, subject.Semester, subject.Stage, subject.Syllabus, subject.Lecturer.ID)

	if err := row.Scan(&subject.ID); err != nil {
		return fmt.Errorf("could not assign id to subject: %v", err)
	}

	return nil
}

func (s *PostgresStore) RemoveSubject(id int) error {
	_, err := s.db.Exec("DELETE FROM subjects WHERE id = $1;", id)
	if err != nil {
		return fmt.Errorf("problem deleting row: %v", err)
	}

	return nil
}

func (s *PostgresStore) UpdateSubject(subject Subject) error {
	if subject.Lecturer == nil {
		subject.Lecturer = &Lecturer{}
	}
	_, err := s.db.Exec(`UPDATE subjects SET name = $1, details = $2, semester = $3, stage = $4, syllabus = $5, lecturer_id = $6 WHERE id = $7;`,
		subject.Name, subject.Details, subject.Semester, subject.Stage, subject.Syllabus, subject.Lecturer.ID, subject.ID)
	if err != nil {
		return fmt.Errorf("problem updating row: %v", err)
	}

	return nil
}

func (s *PostgresStore) GetLecturers() ([]Lecturer, error) {
	rows, err := s.db.Query(`SELECT id, name, image, about, password, subject_id FROM lecturers;`)
	if err != nil {
		return nil, fmt.Errorf("error querying lecturers: %v", err)
	}

	lecturers := []Lecturer{}
	for rows.Next() {
		lecturer := Lecturer{
			Subject: &Subject{},
		}

		err := rows.Scan(&lecturer.ID, &lecturer.Name, &lecturer.Image, &lecturer.About, &lecturer.Password, &lecturer.Subject.ID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning row: %v", err)
		}

		lecturers = append(lecturers, lecturer)
	}
	return lecturers, nil
}

func (s *PostgresStore) GetLecturerByID(id int) (*Lecturer, error) {
	row := s.db.QueryRow("SELECT id, name, image, about, password, subject_id FROM lecturers where id = $1", id)

	lecturer := Lecturer{
		Subject: &Subject{},
	}

	err := row.Scan(&lecturer.ID, &lecturer.Name, &lecturer.Image, &lecturer.About, &lecturer.Password, &lecturer.Subject.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("encountered error while scanning row: %v", err)
	}

	return &lecturer, nil
}

func (s *PostgresStore) AddLecturer(lecturer *Lecturer) error {
	if lecturer.Subject == nil {
		lecturer.Subject = &Subject{}
	}
	row := s.db.QueryRow(`INSERT INTO lecturers (name, image, about, password, subject_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		lecturer.Name, lecturer.Image, lecturer.About, lecturer.Password, lecturer.Subject.ID)

	if err := row.Scan(&lecturer.ID); err != nil {
		return fmt.Errorf("could not assign id to lecturer: %v", err)
	}

	return nil
}

func (s *PostgresStore) RemoveLecturer(id int) error {
	_, err := s.db.Exec("DELETE FROM lecturers WHERE id = $1;", id)
	if err != nil {
		return fmt.Errorf("problem deleting row: %v", err)
	}

	return nil
}

func (s *PostgresStore) UpdateLecturer(lecturer Lecturer) error {
	if lecturer.Subject == nil {
		lecturer.Subject = &Subject{}
	}
	_, err := s.db.Exec(`UPDATE lecturers SET name = $1, image = $2, about = $3, password = $4, subject_id = $5 WHERE id = $6;`,
		lecturer.Name, lecturer.Image, lecturer.About, lecturer.Password, lecturer.Subject.ID, lecturer.ID)
	if err != nil {
		return fmt.Errorf("problem updating row: %v", err)
	}

	return nil
}
