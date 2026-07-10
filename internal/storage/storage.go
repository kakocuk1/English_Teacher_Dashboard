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
		name TEXT NOT NULL,                    -- name of the student, cannot be null
		level TEXT NOT NULL,                   -- CEFR level (e.g., "B2", "C1"), cannot be null
		telegram_id INTEGER NOT NULL DEFAULT 0, -- Telegram user ID, 0 if not linked yet
		lesson_price REAL NOT NULL DEFAULT 0   -- individual price per lesson for this student
	);

	-- schedule table to store lesson information
	CREATE TABLE IF NOT EXISTS lessons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- unique identifier for each lesson
		student_id INTEGER NOT NULL,           -- foreign key referencing the student, cannot be null
		day_of_week TEXT NOT NULL,             -- day of the week (e.g., "Monday", "Tuesday"), cannot be null
		time TEXT NOT NULL,                    -- time of the lesson (e.g., "15:00", "10:00"), cannot be null
		FOREIGN KEY (student_id) REFERENCES students(id)  -- connects to the students table
	);

	-- homework table to store homework information
	CREATE TABLE IF NOT EXISTS homeworks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- unique identifier for each homework
		student_id INTEGER NOT NULL,           -- foreign key referencing the student, cannot be null
		task TEXT NOT NULL,                    -- description of the homework task, cannot be null
		done INTEGER NOT NULL DEFAULT 0,       -- 0 = not done, 1 = done
		created_at TEXT NOT NULL,              -- date when the homework was created (e.g., "26-06-2026")
		FOREIGN KEY (student_id) REFERENCES students(id)
	);
	
	-- lesson_packages stores paid lesson packages for each student
	CREATE TABLE IF NOT EXISTS lesson_packages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- unique identifier for each package
		student_id INTEGER NOT NULL,           -- foreign key referencing the student
		total_lessons INTEGER NOT NULL,        -- total lessons bought (e.g., 4, 8, 12)
		used_lessons INTEGER NOT NULL DEFAULT 0, -- lessons already conducted
		price REAL NOT NULL,                   -- total price paid for the package
		created_at TEXT NOT NULL,              -- date when the package was created
		FOREIGN KEY (student_id) REFERENCES students(id)
	);

	-- lesson_log records each conducted lesson and links it to a package
	CREATE TABLE IF NOT EXISTS lesson_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- unique identifier for each log entry
		student_id INTEGER NOT NULL,           -- foreign key referencing the student
		package_id INTEGER NOT NULL,           -- which package this lesson is taken from
		conducted_at TEXT NOT NULL,            -- date the lesson was conducted
		FOREIGN KEY (student_id) REFERENCES students(id),
		FOREIGN KEY (package_id) REFERENCES lesson_packages(id)
	);`

	_, err := s.db.Exec(query) // execute the query to create tables
	return err                 // return any error that occurs during table creation
}
