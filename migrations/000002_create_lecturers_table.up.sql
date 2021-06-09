CREATE TABLE IF NOT EXISTS lecturers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    image VARCHAR(255),
    password VARCHAR(255),
    about TEXT NOT NULL,
    subject_id int
);

-- ALTER TABLE subjects ADD CONSTRAINT fk_lecturer FOREIGN KEY(lecturer_id) REFERENCES lecturers(id) ON DELETE CASCADE;
-- ALTER TABLE lecturers ADD CONSTRAINT fk_subject FOREIGN KEY(subject_id) REFERENCES subjects(id) ON DELETE CASCADE;
