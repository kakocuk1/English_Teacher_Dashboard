package main

import (
	"log"
	"os"
	"strconv"

	"github.com/kakocuk1/teacher-dashboard/internal/bot"
	"github.com/kakocuk1/teacher-dashboard/internal/service"
	"github.com/kakocuk1/teacher-dashboard/internal/storage"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	teacherID, err := strconv.ParseInt(os.Getenv("TEACHER_TELEGRAM_ID"), 10, 64)
	if err != nil || teacherID == 0 {
		log.Fatal("TEACHER_TELEGRAM_ID is not set or invalid")
	}

	db, err := storage.New(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	svc := service.New(db)

	b, err := bot.New(token, svc, teacherID)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	log.Println("Bot is running...")
	b.Run()
}
