package storage

import "github.com/kakocuk1/teacher-dashboard/internal/model"

// AddStudent adds a new student to the database and returns the created student ID or an error if it fails
func (s *Storage) AddStudent(name, level string) (int, error) {
	result, err := s.db.Exec("INSERT INTO students (name, level) VALUES (?, ?)", name, level)
	if err != nil {
		return 0, err
	}

	// get the last inserted ID to return it
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetStudents retrieves all students from the database and returns a slice of Student or an error if it fails
func (s *Storage) GetStudents() ([]model.Student, error) {
	rows, err := s.db.Query("SELECT id, name, level FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // close the rows after we are done with them

	var students []model.Student

	// iterate through the rows and scan the data into Student structs
	for rows.Next() {
		var st model.Student
		if err := rows.Scan(&st.ID, &st.Name, &st.Level); err != nil {
			return nil, err
		}
		students = append(students, st) // add the student to the slice
	}

	return students, nil
}

// DeleteStudent deletes a student from the database by their ID and returns an error if it fails
func (s *Storage) DeleteStudent(id int) error {
	_, err := s.db.Exec("DELETE FROM students WHERE id = ?", id)
	return err
}
