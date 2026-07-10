package service_test

import (
	"testing"

	"github.com/kakocuk1/teacher-dashboard/internal/service"
)

func TestIsValidLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  bool
	}{
		{name: "valid uppercase", level: "B2", want: true},
		{name: "valid lowercase", level: "a1", want: true},
		{name: "valid with spaces", level: " C1 ", want: true},
		{name: "invalid text", level: "beginner", want: false},
		{name: "empty", level: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.IsValidLevel(tt.level)
			if got != tt.want {
				t.Fatalf("IsValidLevel(%q) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestNormalizeDay(t *testing.T) {
	tests := []struct {
		name string
		day  string
		want string
	}{
		{name: "lowercase", day: "monday", want: "Monday"},
		{name: "uppercase", day: "FRIDAY", want: "Friday"},
		{name: "with spaces", day: " sunday ", want: "Sunday"},
		{name: "empty", day: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.NormalizeDay(tt.day)
			if got != tt.want {
				t.Fatalf("NormalizeDay(%q) = %q, want %q", tt.day, got, tt.want)
			}
		})
	}
}

func TestIsValidDay(t *testing.T) {
	tests := []struct {
		name string
		day  string
		want bool
	}{
		{name: "valid", day: "Monday", want: true},
		{name: "valid lowercase", day: "tuesday", want: true},
		{name: "valid with spaces", day: " friday ", want: true},
		{name: "invalid", day: "Funday", want: false},
		{name: "empty", day: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.IsValidDay(tt.day)
			if got != tt.want {
				t.Fatalf("IsValidDay(%q) = %v, want %v", tt.day, got, tt.want)
			}
		})
	}
}

func TestIsValidLessonTime(t *testing.T) {
	tests := []struct {
		name       string
		lessonTime string
		want       bool
	}{
		{name: "valid", lessonTime: "15:00", want: true},
		{name: "valid with spaces", lessonTime: " 09:30 ", want: true},
		{name: "invalid hour", lessonTime: "25:00", want: false},
		{name: "invalid text", lessonTime: "afternoon", want: false},
		{name: "empty", lessonTime: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.IsValidLessonTime(tt.lessonTime)
			if got != tt.want {
				t.Fatalf("IsValidLessonTime(%q) = %v, want %v", tt.lessonTime, got, tt.want)
			}
		})
	}
}
