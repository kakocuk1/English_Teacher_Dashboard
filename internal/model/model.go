package model

// Student
type Student struct {
	ID    int
	Name  string
	Level string // for example: "beginner", "intermediate", "advanced"
}

// Lesson
type Lesson struct {
	ID        int
	StudentID int
	DayOfWeek string // for example: "Monday", "Tuesday", etc.
	Time      string // for example: "15:00", "10:00", etc.
}

// Homework
type Homework struct {
	ID        int
	StudentID int
	Task      string
	Done      bool
	CreatedAt string // for example: "26-06-2026"
}
