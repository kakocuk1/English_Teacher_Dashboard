package service

import (
	"strings"
	"time"
)

/*
IsValidLevel checks if the student level is one of the common CEFR levels.
This validation belongs to the service layer because both Telegram and future web UI
should follow the same rule.
*/
func IsValidLevel(level string) bool {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "A1", "A2", "B1", "B2", "C1", "C2":
		return true
	default:
		return false
	}
}

/*
NormalizeDay makes the day value easier to validate and store.
For example, "monday" and "MONDAY" both become "Monday".
*/
func NormalizeDay(day string) string {
	day = strings.ToLower(strings.TrimSpace(day))
	if day == "" {
		return ""
	}

	return strings.ToUpper(day[:1]) + day[1:]
}

/*
IsValidDay checks if the lesson day is a valid weekday name.
*/
func IsValidDay(day string) bool {
	switch NormalizeDay(day) {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday":
		return true
	default:
		return false
	}
}

/*
IsValidLessonTime checks the lesson time format.
Go uses the special layout "15:04" to parse time in HH:MM format.
*/
func IsValidLessonTime(lessonTime string) bool {
	_, err := time.Parse("15:04", strings.TrimSpace(lessonTime))
	return err == nil
}
