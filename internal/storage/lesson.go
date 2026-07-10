package storage

import "github.com/kakocuk1/teacher-dashboard/internal/model"

// AddLesson adds a new lesson to the schedule.
func (s *Storage) AddLesson(studentID int, dayOfWeek, lessonTime string) (int, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO lessons (student_id, day_of_week, time) VALUES ($1, $2, $3) RETURNING id",
		studentID, dayOfWeek, lessonTime,
	).Scan(&id)
	return id, err
}

// GetLessonsByStudent returns all lessons for a student.
func (s *Storage) GetLessonsByStudent(studentID int) ([]model.Lesson, error) {
	rows, err := s.db.Query(
		"SELECT id, student_id, day_of_week, time FROM lessons WHERE student_id = $1",
		studentID,
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

// GetAllLessons returns the full schedule sorted by day and time.
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

// DeleteLesson deletes a lesson from the schedule by ID.
func (s *Storage) DeleteLesson(lessonID int) error {
	_, err := s.db.Exec("DELETE FROM lessons WHERE id = $1", lessonID)
	return err
}
