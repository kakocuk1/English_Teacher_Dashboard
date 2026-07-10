package storage

import (
	"time"

	"github.com/kakocuk1/teacher-dashboard/internal/model"
)

// AddHomework adds a new homework to the storage.
func (s *Storage) AddHomework(studentID int, task string) (int, error) {
	result, err := s.db.Exec(
		"INSERT INTO homeworks (student_id, task, created_at) VALUES (?, ?, ?)",
		studentID,
		task,
		time.Now().Format("02-01-2006"),
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

// GetHomeworksByStudent retrieves all homeworks for a given student.
func (s *Storage) GetHomeworksByStudent(studentID int) ([]model.Homework, error) {
	rows, err := s.db.Query(
		"SELECT id, student_id, task, done, created_at FROM homeworks WHERE student_id = ? ORDER BY id DESC",
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return homeworks, nil
}

// MarkHomeworkDone marks a homework as done.
func (s *Storage) MarkHomeworkDone(homeworkID int) error {
	_, err := s.db.Exec(
		"UPDATE homeworks SET done = 1 WHERE id = ?",
		homeworkID,
	)
	return err
}

// DeleteHomework deletes a homework from the storage.
func (s *Storage) DeleteHomework(homeworkID int) error {
	_, err := s.db.Exec(
		"DELETE FROM homeworks WHERE id = ?",
		homeworkID,
	)
	return err
}
