package model

// Student represents a student in the system.
type Student struct {
	ID          int
	Name        string
	Level       string  // CEFR level: A1, A2, B1, B2, C1, C2
	TelegramID  int64   // Telegram user ID, 0 if not linked yet
	LessonPrice float64 // individual price per lesson
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

// LessonPackage represents a paid package of lessons for a student.
type LessonPackage struct {
	ID           int
	StudentID    int
	TotalLessons int     // total lessons bought e.g. 4, 8, 12
	UsedLessons  int     // lessons already conducted
	Price        float64 // total price paid for the package
	CreatedAt    string
}

// LessonLog represents a single conducted lesson.
type LessonLog struct {
	ID          int
	StudentID   int
	PackageID   int    // which package this lesson is taken from
	ConductedAt string // date the lesson was conducted
}
