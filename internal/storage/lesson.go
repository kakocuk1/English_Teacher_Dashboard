package storage

import "github.com/kakocuk1/teacher-dashboard/internal/model"

// AddLesson adds a new lesson to the storage.
func (s *Storage) AddLesson(studentID int, dayOfWeek, lessonTime string) (int, error) {
	result, err := s.db.Exec(
		"INSERT INTO lessons (student_id, day_of_week, time) VALUES (?, ?, ?)",
		studentID,
		dayOfWeek,
		lessonTime,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetLessonsByStudent retrieves all lessons for a given student.
func (s *Storage) GetLessonsByStudent(studentID int) ([]model.Lesson, error) {
	rows, err := s.db.Query(
		"SELECT id, student_id, day_of_week, time FROM lessons WHERE student_id = ?",
		studentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []model.Lesson

	for rows.Next() {
		var l model.Lesson // l=Lesson
		if err := rows.Scan(&l.ID, &l.StudentID, &l.DayOfWeek, &l.Time); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}

	return lessons, nil
}

// GetAllLessons retrieves all lessons from the storage.
func (s *Storage) GetAllLessons() ([]model.Lesson, error) {
	rows, err := s.db.Query(
		"SELECT id, student_id, day_of_week, time FROM lessons ORDER BY day_of_week, time",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []model.Lesson

	for rows.Next() {
		var l model.Lesson
		if err := rows.Scan(&l.ID, &l.StudentID, &l.DayOfWeek, &l.Time); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

// DeleteLesson deletes a lesson by its ID.
func (s *Storage) DeleteLesson(lessonID int) error {
	_, err := s.db.Exec(
		"DELETE FROM lessons WHERE id = ?",
		lessonID,
	)
	return err
}
