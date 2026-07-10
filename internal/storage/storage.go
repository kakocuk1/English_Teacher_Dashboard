package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Storage struct to hold the database connection
type Storage struct {
	db *sql.DB
}

// New create a new connection to the database and return a Storage or error if it fails
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path) // open the database connection
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil { // check if the connection is alive
		return nil, err
	}

	s := &Storage{db: db} // create a new Storage instance with the database connection

	if err := s.migrate(); err != nil { // run the migration to create tables if they don't exist
		return nil, err
	}

	return s, nil // return the Storage instance
}

// Migrate creates the necessary tables in the database if they do not exist
func (s *Storage) migrate() error {
	query := `
	-- Table of students to store student information
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- unique identifier for each student
		name TEXT NOT NULL,  -- name of the student, cannot be null
		level TEXT NOT NULL  -- level of the student (e.g., "beginner", "intermediate", "advanced"), cannot be null
	);

	-- schedule table to store lesson information
	CREATE TABLE IF NOT EXISTS lessons (
	id INTEGER PRIMARY KEY AUTOINCREMENT,  -- unique identifier for each lesson
	student_id INTEGER NOT NULL,  -- foreign key referencing the student who has the lesson, cannot be null
	day_of_week TEXT NOT NULL,  -- day of the week for the lesson (e.g., "Monday", "Tuesday"), cannot be null
	time TEXT NOT NULL,  -- time of the lesson (e.g., "15:00", "10:00"), cannot be null
	FOREIGN KEY (student_id) REFERENCES students(id)  --connects the student_id in lessons to the id in students
	);

	-- homework table to store homework information
	CREATE TABLE IF NOT EXISTS homeworks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,  -- unique identifier for each homework
	student_id INTEGER NOT NULL,  -- foreign key referencing the student who has the homework, cannot be null
	task TEXT NOT NULL,  -- description of the homework task, cannot be null
	done INTEGER NOT NULL DEFAULT 0,  -- indicates whether the homework is done (0 for false, 1 for true), cannot be null
	created_at TEXT NOT NULL,  -- date when the homework was created (e.g., "26-06-2026"), cannot be null
	FOREIGN KEY (student_id) REFERENCES students(id)  -- connects the student_id in homeworks to the id in students
	);
	`

	_, err := s.db.Exec(query) // execute the query to create tables
	return err                 // return any error that occurs during table creation
}
