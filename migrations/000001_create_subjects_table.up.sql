CREATE TABLE IF NOT EXISTS lecturers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    image VARCHAR(255),
    password VARCHAR(255),
    about TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS subjects(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    details TEXT NOT NULL,
    syllabus TEXT NOT NULL,
    semester INT default 1,
    stage INT default 1,
    lecturer_id INT NOT NULL,
    FOREIGN KEY (lecturer_id) REFERENCES lecturers (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS announcements(
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    details VARCHAR(100),
    image VARCHAR(200)
);
