CREATE TABLE IF NOT EXISTS subjects(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    details TEXT NOT NULL,
    syllabus TEXT NOT NULL,
    semester INT,
    stage INT,
    lecturer_id INT
);
