package main

import (
	"log"
	"os"

	"github.com/kakocuk1/teacher-dashboard/internal/bot"
	"github.com/kakocuk1/teacher-dashboard/internal/service"
	"github.com/kakocuk1/teacher-dashboard/internal/storage"
)

func main() {
	// take a token from environment variable, DO NOT SAVE HERE TOKEN
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set")
	}

	// Create connection to the database
	db, err := storage.New("teacher.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create a service with business logic
	svc := service.New(db)

	// Create and start the bot
	b, err := bot.New(token, svc)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Println("Bot is running...")
	b.Run()
}
