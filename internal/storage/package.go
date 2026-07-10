package storage

import (
	"time"

	"github.com/kakocuk1/teacher-dashboard/internal/model"
)

// AddLessonPackage adds a new paid lesson package for a student.
func (s *Storage) AddLessonPackage(studentID, totalLessons int, price float64) (int, error) {
	result, err := s.db.Exec(
		`INSERT INTO lesson_packages (student_id, total_lessons, used_lessons, price, created_at)
		VALUES (?, ?, 0, ?, ?)`,
		studentID,
		totalLessons,
		price,
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

// GetActivePackage returns the current active package for a student.
// Active means it still has lessons remaining.
func (s *Storage) GetActivePackage(studentID int) (*model.LessonPackage, error) {
	row := s.db.QueryRow(
		`SELECT id, student_id, total_lessons, used_lessons, price, created_at
		FROM lesson_packages
		WHERE student_id = ? AND used_lessons < total_lessons
		ORDER BY created_at ASC
		LIMIT 1`,
		studentID,
	)

	var p model.LessonPackage
	err := row.Scan(&p.ID, &p.StudentID, &p.TotalLessons, &p.UsedLessons, &p.Price, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetAllPackages returns all packages for a student including finished ones.
func (s *Storage) GetAllPackages(studentID int) ([]model.LessonPackage, error) {
	rows, err := s.db.Query(
		`SELECT id, student_id, total_lessons, used_lessons, price, created_at
		FROM lesson_packages
		WHERE student_id = ?
		ORDER BY created_at DESC`,
		studentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []model.LessonPackage
	for rows.Next() {
		var p model.LessonPackage
		if err := rows.Scan(&p.ID, &p.StudentID, &p.TotalLessons, &p.UsedLessons, &p.Price, &p.CreatedAt); err != nil {
			return nil, err
		}
		packages = append(packages, p)
	}

	return packages, nil
}

// MarkLessonConducted records a conducted lesson and increments used_lessons in the active package.
func (s *Storage) MarkLessonConducted(studentID, packageID int) error {
	tx, err := s.db.Begin() // start a transaction so both queries succeed or fail together
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO lesson_log (student_id, package_id, conducted_at) VALUES (?, ?, ?)`,
		studentID,
		packageID,
		time.Now().Format("02-01-2006"),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		`UPDATE lesson_packages SET used_lessons = used_lessons + 1 WHERE id = ?`,
		packageID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
