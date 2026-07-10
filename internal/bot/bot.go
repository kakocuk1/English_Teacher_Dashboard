package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kakocuk1/teacher-dashboard/internal/service"
)

// Bot represents the Telegram bot structure.
type Bot struct {
	api         *tgbotapi.BotAPI
	service     *service.Service
	pendingLink int // student ID waiting to be linked to the next /start from a student
}

// New creates a new instance of the Bot.
func New(token string, service *service.Service) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:     api,
		service: service,
	}, nil
}

// Run starts the bot and listens for incoming messages.
func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			log.Printf("ignoring non-command message from chat %d", update.Message.Chat.ID)
			continue
		}

		b.handleMessage(update.Message)
	}
}
