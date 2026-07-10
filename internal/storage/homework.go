package storage

import (
	"time"

	"github.com/kakocuk1/teacher-dashboard/internal/model"
)

// AddHomework adds a new homework task for a student.
func (s *Storage) AddHomework(studentID int, task string) (int, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO homeworks (student_id, task, done, created_at) VALUES ($1, $2, FALSE, $3) RETURNING id",
		studentID,
		task,
		time.Now().Format("02-01-2006"),
	).Scan(&id)
	return id, err
}

// GetHomeworksByStudent returns all homework tasks for a student.
func (s *Storage) GetHomeworksByStudent(studentID int) ([]model.Homework, error) {
	rows, err := s.db.Query(
		"SELECT id, student_id, task, done, created_at FROM homeworks WHERE student_id = $1",
		studentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var homeworks []model.Homework
	for rows.Next() {
		var hw model.Homework
		if err := rows.Scan(&hw.ID, &hw.StudentID, &hw.Task, &hw.Done, &hw.CreatedAt); err != nil {
			return nil, err
		}
		homeworks = append(homeworks, hw)
	}

	return homeworks, nil
}

// MarkHomeworkDone marks a homework task as done.
func (s *Storage) MarkHomeworkDone(homeworkID int) error {
	_, err := s.db.Exec(
		"UPDATE homeworks SET done = TRUE WHERE id = $1",
		homeworkID,
	)
	return err
}

// DeleteHomework deletes a homework task by ID.
func (s *Storage) DeleteHomework(homeworkID int) error {
	_, err := s.db.Exec("DELETE FROM homeworks WHERE id = $1", homeworkID)
	return err
}
