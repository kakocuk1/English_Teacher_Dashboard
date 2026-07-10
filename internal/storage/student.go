package storage

import "github.com/kakocuk1/teacher-dashboard/internal/model"

// AddStudent adds a new student to the database and returns the created student ID or an error if it fails
func (s *Storage) AddStudent(name, level string) (int, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO students (name, level) VALUES ($1, $2) RETURNING id",
		name, level,
	).Scan(&id)
	return id, err
}

// GetStudents retrieves all students from the database and returns a slice of Student or an error if it fails
func (s *Storage) GetStudents() ([]model.Student, error) {
	rows, err := s.db.Query("SELECT id, name, level, telegram_id, lesson_price FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var st model.Student
		if err := rows.Scan(&st.ID, &st.Name, &st.Level, &st.TelegramID, &st.LessonPrice); err != nil {
			return nil, err
		}
		students = append(students, st)
	}

	return students, nil
}

// DeleteStudent deletes a student from the database by their ID and returns an error if it fails
func (s *Storage) DeleteStudent(id int) error {
	_, err := s.db.Exec("DELETE FROM students WHERE id = $1", id)
	return err
}

// GetStudentByID returns a student by their ID.
func (s *Storage) GetStudentByID(id int) (*model.Student, error) {
	row := s.db.QueryRow(
		"SELECT id, name, level, telegram_id, lesson_price FROM students WHERE id = $1",
		id,
	)

	var st model.Student
	err := row.Scan(&st.ID, &st.Name, &st.Level, &st.TelegramID, &st.LessonPrice)
	if err != nil {
		return nil, err
	}

	return &st, nil
}
