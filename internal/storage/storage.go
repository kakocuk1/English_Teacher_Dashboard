package storage

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Storage holds the database connection.
type Storage struct {
	db *sql.DB
}

// New creates a new connection to the PostgreSQL database.
// dsn is the connection string, e.g. "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
func New(dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Storage{db: db}

	if err := s.migrate(); err != nil {
		return nil, err
	}

	return s, nil
}

// migrate creates the necessary tables if they do not exist.
func (s *Storage) migrate() error {
	query := `
	-- Table of students to store student information
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,                  -- auto-incrementing unique identifier
		name TEXT NOT NULL,                     -- name of the student
		level TEXT NOT NULL,                    -- CEFR level (e.g., "B2", "C1")
		telegram_id BIGINT NOT NULL DEFAULT 0,  -- Telegram user ID, 0 if not linked yet
		lesson_price NUMERIC NOT NULL DEFAULT 0 -- individual price per lesson
	);

	-- schedule table to store lesson information
	CREATE TABLE IF NOT EXISTS lessons (
		id SERIAL PRIMARY KEY,
		student_id INTEGER NOT NULL,
		day_of_week TEXT NOT NULL,              -- day of the week (e.g., "Monday")
		time TEXT NOT NULL,                     -- time of the lesson (e.g., "15:00")
		FOREIGN KEY (student_id) REFERENCES students(id)
	);

	-- homework table to store homework information
	CREATE TABLE IF NOT EXISTS homeworks (
		id SERIAL PRIMARY KEY,
		student_id INTEGER NOT NULL,
		task TEXT NOT NULL,                     -- description of the homework task
		done BOOLEAN NOT NULL DEFAULT FALSE,    -- false = not done, true = done
		created_at TEXT NOT NULL,
		FOREIGN KEY (student_id) REFERENCES students(id)
	);

	-- lesson_packages stores paid lesson packages for each student
	CREATE TABLE IF NOT EXISTS lesson_packages (
		id SERIAL PRIMARY KEY,
		student_id INTEGER NOT NULL,
		total_lessons INTEGER NOT NULL,         -- total lessons bought (e.g., 4, 8, 12)
		used_lessons INTEGER NOT NULL DEFAULT 0,
		price NUMERIC NOT NULL,                 -- total price paid for the package
		created_at TEXT NOT NULL,
		FOREIGN KEY (student_id) REFERENCES students(id)
	);

	-- lesson_log records each conducted lesson
	CREATE TABLE IF NOT EXISTS lesson_log (
		id SERIAL PRIMARY KEY,
		student_id INTEGER NOT NULL,
		package_id INTEGER NOT NULL,            -- which package this lesson is taken from
		conducted_at TEXT NOT NULL,
		FOREIGN KEY (student_id) REFERENCES students(id),
		FOREIGN KEY (package_id) REFERENCES lesson_packages(id)
	);`

	_, err := s.db.Exec(query)
	return err
}
