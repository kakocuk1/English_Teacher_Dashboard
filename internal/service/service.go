package service

import (
	"fmt"
	"strings"

	"github.com/kakocuk1/teacher-dashboard/internal/model"
)

// Business logic layer

// Storage is an interface that defines the methods for interacting with the data storage.
// The service uses this interface so different interfaces, such as Telegram bot or web UI,
// can share the same business logic without knowing how the database works.
type Storage interface {
	// Methods for managing students.
	AddStudent(name, level string) (int, error)
	GetStudents() ([]model.Student, error)
	DeleteStudent(id int) error
	LinkStudent(studentID int, telegramID int64) error
	GetStudentByTelegramID(telegramID int64) (*model.Student, error)
	SetLessonPrice(studentID int, price float64) error

	// Methods for homework management.
	AddHomework(studentID int, task string) (int, error)
	GetHomeworksByStudent(studentID int) ([]model.Homework, error)
	MarkHomeworkDone(homeworkID int) error
	DeleteHomework(homeworkID int) error

	// Methods for scheduling lessons.
	AddLesson(studentID int, dayOfWeek, lessonTime string) (int, error)
	GetLessonsByStudent(studentID int) ([]model.Lesson, error)
	GetAllLessons() ([]model.Lesson, error)
	DeleteLesson(lessonID int) error

	// Methods for lesson packages and payment tracking.
	AddLessonPackage(studentID, totalLessons int, price float64) (int, error)
	GetActivePackage(studentID int) (*model.LessonPackage, error)
	GetAllPackages(studentID int) ([]model.LessonPackage, error)
	MarkLessonConducted(studentID, packageID int) error
}

// Service contains business rules of the application.
// Telegram handlers and future web handlers should call this layer instead of storage directly.
type Service struct {
	storage Storage
}

// New creates a new Service instance.
func New(storage Storage) *Service {
	return &Service{storage: storage}
}

// AddStudent validates and adds a new student to the system.
func (s *Service) AddStudent(name, level string) (int, error) {
	name = strings.TrimSpace(name)
	level = strings.ToUpper(strings.TrimSpace(level))

	if name == "" {
		return 0, fmt.Errorf("student name cannot be empty")
	}
	if !IsValidLevel(level) {
		return 0, fmt.Errorf("invalid student level: %s", level)
	}

	return s.storage.AddStudent(name, level)
}

// GetStudents returns all students from the system.
func (s *Service) GetStudents() ([]model.Student, error) {
	return s.storage.GetStudents()
}

// DeleteStudent deletes a student by ID.
func (s *Service) DeleteStudent(id int) error {
	if id <= 0 {
		return fmt.Errorf("student ID must be greater than zero")
	}

	return s.storage.DeleteStudent(id)
}

// AddHomework validates and adds homework to a student.
func (s *Service) AddHomework(studentID int, task string) (int, error) {
	task = strings.TrimSpace(task)

	if studentID <= 0 {
		return 0, fmt.Errorf("student ID must be greater than zero")
	}
	if task == "" {
		return 0, fmt.Errorf("homework task cannot be empty")
	}

	return s.storage.AddHomework(studentID, task)
}

// GetHomeworksByStudent returns all homework tasks for one student.
func (s *Service) GetHomeworksByStudent(studentID int) ([]model.Homework, error) {
	if studentID <= 0 {
		return nil, fmt.Errorf("student ID must be greater than zero")
	}

	return s.storage.GetHomeworksByStudent(studentID)
}

// MarkHomeworkDone marks a homework task as done.
func (s *Service) MarkHomeworkDone(homeworkID int) error {
	if homeworkID <= 0 {
		return fmt.Errorf("homework ID must be greater than zero")
	}

	return s.storage.MarkHomeworkDone(homeworkID)
}

// DeleteHomework deletes a homework task.
func (s *Service) DeleteHomework(homeworkID int) error {
	if homeworkID <= 0 {
		return fmt.Errorf("homework ID must be greater than zero")
	}

	return s.storage.DeleteHomework(homeworkID)
}

// AddLesson validates and adds a lesson to the schedule.
func (s *Service) AddLesson(studentID int, dayOfWeek, lessonTime string) (int, error) {
	dayOfWeek = NormalizeDay(dayOfWeek)
	lessonTime = strings.TrimSpace(lessonTime)

	if studentID <= 0 {
		return 0, fmt.Errorf("student ID must be greater than zero")
	}
	if !IsValidDay(dayOfWeek) {
		return 0, fmt.Errorf("invalid lesson day: %s", dayOfWeek)
	}
	if !IsValidLessonTime(lessonTime) {
		return 0, fmt.Errorf("invalid lesson time: %s", lessonTime)
	}

	return s.storage.AddLesson(studentID, dayOfWeek, lessonTime)
}

// GetLessonsByStudent returns lessons for one student.
func (s *Service) GetLessonsByStudent(studentID int) ([]model.Lesson, error) {
	if studentID <= 0 {
		return nil, fmt.Errorf("student ID must be greater than zero")
	}

	return s.storage.GetLessonsByStudent(studentID)
}

// GetAllLessons returns the schedule of all students.
func (s *Service) GetAllLessons() ([]model.Lesson, error) {
	return s.storage.GetAllLessons()
}

// DeleteLesson deletes a lesson from the schedule.
func (s *Service) DeleteLesson(lessonID int) error {
	if lessonID <= 0 {
		return fmt.Errorf("lesson ID must be greater than zero")
	}

	return s.storage.DeleteLesson(lessonID)
}

// LinkStudent links a Telegram user ID to a student.
func (s *Service) LinkStudent(studentID int, telegramID int64) error {
	if studentID <= 0 {
		return fmt.Errorf("student ID must be greater than zero")
	}
	if telegramID == 0 {
		return fmt.Errorf("telegram ID cannot be zero")
	}
	return s.storage.LinkStudent(studentID, telegramID)
}

// GetStudentByTelegramID finds a student by their Telegram user ID.
func (s *Service) GetStudentByTelegramID(telegramID int64) (*model.Student, error) {
	return s.storage.GetStudentByTelegramID(telegramID)
}

// SetLessonPrice sets the individual lesson price for a student.
func (s *Service) SetLessonPrice(studentID int, price float64) error {
	if studentID <= 0 {
		return fmt.Errorf("student ID must be greater than zero")
	}
	if price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	return s.storage.SetLessonPrice(studentID, price)
}

// AddLessonPackage adds a new paid lesson package for a student.
func (s *Service) AddLessonPackage(studentID, totalLessons int, price float64) (int, error) {
	if studentID <= 0 {
		return 0, fmt.Errorf("student ID must be greater than zero")
	}
	if totalLessons <= 0 {
		return 0, fmt.Errorf("total lessons must be greater than zero")
	}
	if price < 0 {
		return 0, fmt.Errorf("price cannot be negative")
	}
	return s.storage.AddLessonPackage(studentID, totalLessons, price)
}

// GetActivePackage returns the current active lesson package for a student.
func (s *Service) GetActivePackage(studentID int) (*model.LessonPackage, error) {
	if studentID <= 0 {
		return nil, fmt.Errorf("student ID must be greater than zero")
	}
	return s.storage.GetActivePackage(studentID)
}

// GetAllPackages returns all lesson packages for a student.
func (s *Service) GetAllPackages(studentID int) ([]model.LessonPackage, error) {
	if studentID <= 0 {
		return nil, fmt.Errorf("student ID must be greater than zero")
	}
	return s.storage.GetAllPackages(studentID)
}

/*
ConductLesson marks a lesson as conducted for a student.
It finds the active package and records the lesson in the log.
Returns true if only 1 lesson remains after this one — signals to remind about payment.
*/
func (s *Service) ConductLesson(studentID int) (remindPayment bool, err error) {
	if studentID <= 0 {
		return false, fmt.Errorf("student ID must be greater than zero")
	}

	pkg, err := s.storage.GetActivePackage(studentID)
	if err != nil {
		return false, fmt.Errorf("no active package for this student")
	}

	if err := s.storage.MarkLessonConducted(studentID, pkg.ID); err != nil {
		return false, err
	}

	// check if only 1 lesson remains after marking this one
	remaining := pkg.TotalLessons - pkg.UsedLessons - 1
	return remaining == 1, nil
}
