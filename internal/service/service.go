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
