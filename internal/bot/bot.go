package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kakocuk1/teacher-dashboard/internal/service"
)

// Bot represents the Telegram bot structure.
type Bot struct {
	api         *tgbotapi.BotAPI
	service     *service.Service
	pendingLink int   // student ID waiting to be linked to the next /start from a student
	teacherID   int64 // only this user can use teacher commands
}

// New creates a new instance of the Bot.
func New(token string, service *service.Service, teacherID int64) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:       api,
		service:   service,
		teacherID: teacherID,
	}, nil
}

// Run starts the bot and listens for incoming messages.
func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			b.handleMessage(update.Message)
			continue
		}

		// handle button clicks
		if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

// isTeacher checks if the message is from the teacher.
func (b *Bot) isTeacher(message *tgbotapi.Message) bool {
	return message.From.ID == b.teacherID
}
