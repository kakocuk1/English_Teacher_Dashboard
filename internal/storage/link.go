package storage

import "github.com/kakocuk1/teacher-dashboard/internal/model"

// LinkStudent links a Telegram user ID to a student by student ID.
func (s *Storage) LinkStudent(studentID int, telegramID int64) error {
	_, err := s.db.Exec(
		`UPDATE students SET telegram_id = ? WHERE id = ?`,
		telegramID,
		studentID,
	)
	return err
}

// GetStudentByTelegramID finds a student by their Telegram user ID.
// Used when a student sends a message to the bot.
func (s *Storage) GetStudentByTelegramID(telegramID int64) (*model.Student, error) {
	row := s.db.QueryRow(
		`SELECT id, name, level, telegram_id, lesson_price FROM students WHERE telegram_id = ?`,
		telegramID,
	)

	var st model.Student
	err := row.Scan(&st.ID, &st.Name, &st.Level, &st.TelegramID, &st.LessonPrice)
	if err != nil {
		return nil, err
	}

	return &st, nil
}

// SetLessonPrice sets the individual lesson price for a student.
func (s *Storage) SetLessonPrice(studentID int, price float64) error {
	_, err := s.db.Exec(
		`UPDATE students SET lesson_price = ? WHERE id = ?`,
		price,
		studentID,
	)
	return err
}
