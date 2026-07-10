package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kakocuk1/teacher-dashboard/internal/service"
)

const (
	emojiCheck    = "\u2705"
	emojiCross    = "\u274c"
	emojiBook     = "\U0001f4da"
	emojiStudent  = "\U0001f468\u200d\U0001f393"
	emojiCalendar = "\U0001f4c5"
	emojiWarning  = "\u26a0\ufe0f"
	emojiSuccess  = "\U0001f389"
	emojiError    = "\U0001f6ab"
	emojiInfo     = "\u2139\ufe0f"
)

/*
handleMessage is the main router for Telegram commands.
It reads the command from the incoming message and calls the correct handler.
If the command is unknown, it sends a short help message to the user.
*/
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.handleStart(message)
	case "my_homeworks", "my_lessons":
		b.handleStudentCommands(message)
	default:
		if !b.isTeacher(message) {
			b.send(message.Chat.ID, fmt.Sprintf("%s You are not authorized to use this command.", emojiError))
			return
		}
		b.handleTeacherCommands(message)
	}
}

func (b *Bot) handleStudentCommands(message *tgbotapi.Message) {
	switch message.Command() {
	case "my_homeworks":
		b.handleMyHomeworks(message)
	case "my_lessons":
		b.handleMyLessons(message)
	}
}

func (b *Bot) handleTeacherCommands(message *tgbotapi.Message) {
	switch message.Command() {
	case "add_student":
		b.handleAddStudent(message)
	case "students":
		b.handleGetStudents(message)
	case "add_homework":
		b.handleAddHomework(message)
	case "homeworks":
		b.handleGetHomeworks(message)
	case "done":
		b.handleMarkDone(message)
	case "schedule":
		b.handleGetSchedule(message)
	case "add_lesson":
		b.handleAddLesson(message)
	case "link_student":
		b.handleLinkStudent(message)
	case "set_price":
		b.handleSetPrice(message)
	case "add_package":
		b.handleAddPackage(message)
	case "lesson_done":
		b.handleLessonDone(message)
	case "balance":
		b.handleBalance(message)
	default:
		b.send(message.Chat.ID, fmt.Sprintf("%s Unknown command. Use /start to see all commands.", emojiWarning))
	}
}

/*
send is a small helper for sending text messages.
It creates a Telegram message object and sends it through the bot API.
If Telegram returns an error, the error is written to the application log.
*/
func (b *Bot) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("error sending message: %v", err)
	}
}

// sendWithKeyboard sends a message with inline keyboard buttons.
func (b *Bot) sendWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("error sending message with keyboard: %v", err)
	}
}

// buildStudentKeyboard builds an inline keyboard with all students as buttons.
// action is a prefix for callback data, e.g. "lesson_done" or "balance"
func (b *Bot) buildStudentKeyboard(action string) (tgbotapi.InlineKeyboardMarkup, error) {
	students, err := b.service.GetStudents()
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, s := range students {
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s (%s)", s.Name, s.Level),
			fmt.Sprintf("%s:%d", action, s.ID), // callback data: "lesson_done:2"
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}

/*
handleStart sends the welcome message and the list of available commands.
This is usually the first command a user sends to the bot: /start.
*/
func (b *Bot) handleStart(message *tgbotapi.Message) {
	// pending link — linking a student account
	if b.pendingLink != 0 && !b.isTeacher(message) {
		if err := b.service.LinkStudent(b.pendingLink, message.From.ID); err != nil {
			b.send(message.Chat.ID, fmt.Sprintf("%s Error linking account: %s", emojiError, err.Error()))
			return
		}
		b.send(message.Chat.ID, fmt.Sprintf("%s Your account has been linked! Use /my_homeworks and /my_lessons.", emojiSuccess))
		b.pendingLink = 0
		return
	}

	// teacher menu
	if b.isTeacher(message) {
		text := fmt.Sprintf(`%s Teacher Dashboard

/students — list of all students
/add_student Name Level — add a student
/link_student StudentID — link next student who writes /start
/set_price StudentID Price — set lesson price
/add_package StudentID Lessons Price — add paid package
/lesson_done StudentID — mark lesson as conducted
/balance StudentID — show remaining lessons
/add_homework StudentID Task — assign homework
/homeworks StudentID — view student homework
/done HomeworkID — mark homework as done
/add_lesson StudentID Day Time — add lesson to schedule
/schedule — view full schedule`, emojiStudent)
		b.send(message.Chat.ID, text)
		return
	}

	// linked student menu
	student, err := b.service.GetStudentByTelegramID(message.From.ID)
	if err == nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Welcome, %s!\n\n/my_homeworks — your homework\n/my_lessons — remaining lessons", emojiStudent, student.Name))
		return
	}

	// unknown user
	b.send(message.Chat.ID, fmt.Sprintf("%s You are not authorized.", emojiError))
}

/*
handleAddStudent creates a new student from the /add_student command.
The expected format is: /add_student Name Level.
The last word is used as the level, and everything before it is used as the name.
Example: /add_student John Doe B2 -> name = "John Doe", level = "B2".
*/
func (b *Bot) handleAddStudent(message *tgbotapi.Message) {
	args := strings.Fields(message.CommandArguments())
	if len(args) < 2 {
		b.send(message.Chat.ID, fmt.Sprintf("%s Usage: /add_student Name Level\nExample: /add_student John Doe B2", emojiWarning))
		return
	}

	level := strings.ToUpper(args[len(args)-1])
	name := strings.Join(args[:len(args)-1], " ")

	if !service.IsValidLevel(level) {
		b.send(message.Chat.ID, fmt.Sprintf("%s Invalid level. Use one of: A1, A2, B1, B2, C1, C2", emojiError))
		return
	}

	id, err := b.service.AddStudent(name, level)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error adding student: %s", emojiError, err.Error()))
		return
	}

	b.send(message.Chat.ID, fmt.Sprintf("%s Student %s (%s) added with ID %d", emojiSuccess, name, level, id))
}

/*
handleGetStudents shows all students saved in the system.
It asks the service layer for the student list and formats the result as text.
If there are no students yet, it tells the user how to add one.
*/
func (b *Bot) handleGetStudents(message *tgbotapi.Message) {
	students, err := b.service.GetStudents()
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error getting students: %s", emojiError, err.Error()))
		return
	}

	if len(students) == 0 {
		b.send(message.Chat.ID, fmt.Sprintf("%s No students yet. Add one with /add_student", emojiInfo))
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s Students:\n\n", emojiStudent))
	for _, s := range students {
		sb.WriteString(fmt.Sprintf("ID %d - %s (%s)\n", s.ID, s.Name, s.Level))
	}

	b.send(message.Chat.ID, sb.String())
}

/*
handleAddHomework creates a new homework task for a student.
The expected format is: /add_homework StudentID Task.
SplitN is used because the task can contain spaces.
Example: /add_homework 1 Read unit 5.
*/
func (b *Bot) handleAddHomework(message *tgbotapi.Message) {
	args := strings.SplitN(message.CommandArguments(), " ", 2)
	if len(args) < 2 {
		b.send(message.Chat.ID, fmt.Sprintf("%s Usage: /add_homework StudentID Task\nExample: /add_homework 1 Read unit 5", emojiWarning))
		return
	}

	studentID, err := strconv.Atoi(args[0])
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s StudentID must be a number", emojiError))
		return
	}

	task := strings.TrimSpace(args[1])
	if task == "" {
		b.send(message.Chat.ID, fmt.Sprintf("%s Homework task cannot be empty", emojiError))
		return
	}

	id, err := b.service.AddHomework(studentID, task)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error adding homework: %s", emojiError, err.Error()))
		return
	}

	b.send(message.Chat.ID, fmt.Sprintf("%s Homework added with ID %d", emojiSuccess, id))

	// notify the student if their Telegram account is linked
	student, err := b.service.GetStudentByID(studentID)
	if err == nil && student.TelegramID != 0 {
		b.send(student.TelegramID, fmt.Sprintf("%s New homework from your teacher:\n\n%s", emojiBook, task))
	}
}

/*
handleGetHomeworks shows all homework tasks for one student.
The command must contain a student ID: /homeworks StudentID.
Each homework item is printed with its ID, status, task text, and creation date.
*/
func (b *Bot) handleGetHomeworks(message *tgbotapi.Message) {
	keyboard, err := b.buildStudentKeyboard("homeworks")
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error getting students: %s", emojiError, err.Error()))
		return
	}

	b.sendWithKeyboard(message.Chat.ID, "Select a student:", keyboard)
}

/*
handleMarkDone marks one homework task as completed.
The command must contain a homework ID: /done HomeworkID.
The real update is done in the service layer.
*/
func (b *Bot) handleMarkDone(message *tgbotapi.Message) {
	homeworkID, err := strconv.Atoi(strings.TrimSpace(message.CommandArguments()))
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Usage: /done HomeworkID\nExample: /done 1", emojiWarning))
		return
	}

	if err := b.service.MarkHomeworkDone(homeworkID); err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error marking homework as done: %s", emojiError, err.Error()))
		return
	}

	b.send(message.Chat.ID, fmt.Sprintf("%s Homework marked as done!", emojiCheck))
}

/*
handleGetSchedule shows the full lesson schedule.
It gets all lessons from the service layer and formats them for Telegram.
If there are no lessons yet, it tells the user how to add one.
*/
func (b *Bot) handleGetSchedule(message *tgbotapi.Message) {
	lessons, err := b.service.GetAllLessons()
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error getting schedule: %s", emojiError, err.Error()))
		return
	}

	if len(lessons) == 0 {
		b.send(message.Chat.ID, fmt.Sprintf("%s No lessons yet. Add one with /add_lesson", emojiInfo))
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s Schedule:\n\n", emojiCalendar))
	for _, l := range lessons {
		sb.WriteString(fmt.Sprintf("%s %s - Student ID %d\n", l.DayOfWeek, l.Time, l.StudentID))
	}

	b.send(message.Chat.ID, sb.String())
}

/*
handleAddLesson adds a new lesson to the schedule.
The expected format is: /add_lesson StudentID Day Time.
Example: /add_lesson 1 Monday 15:00.
*/
func (b *Bot) handleAddLesson(message *tgbotapi.Message) {
	args := strings.Fields(message.CommandArguments())
	if len(args) != 3 {
		b.send(message.Chat.ID, fmt.Sprintf("%s Usage: /add_lesson StudentID Day Time\nExample: /add_lesson 1 Monday 15:00", emojiWarning))
		return
	}

	studentID, err := strconv.Atoi(args[0])
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s StudentID must be a number", emojiError))
		return
	}

	day := service.NormalizeDay(args[1])
	lessonTime := args[2]

	if !service.IsValidDay(day) {
		b.send(message.Chat.ID, fmt.Sprintf("%s Invalid day. Use: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday", emojiError))
		return
	}

	if !service.IsValidLessonTime(lessonTime) {
		b.send(message.Chat.ID, fmt.Sprintf("%s Invalid time format. Use HH:MM, for example 15:00", emojiError))
		return
	}

	id, err := b.service.AddLesson(studentID, day, lessonTime)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error adding lesson: %s", emojiError, err.Error()))
		return
	}

	b.send(message.Chat.ID, fmt.Sprintf("%s Lesson added with ID %d", emojiSuccess, id))
}

// handleLinkStudent links the next student who writes /start to a student ID.
// Usage: /link_student StudentID
func (b *Bot) handleLinkStudent(message *tgbotapi.Message) {
	studentID, err := strconv.Atoi(strings.TrimSpace(message.CommandArguments()))
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Usage: /link_student StudentID\nExample: /link_student 1", emojiWarning))
		return
	}

	b.pendingLink = studentID // remember which student ID to link on next /start
	b.send(message.Chat.ID, fmt.Sprintf("%s Ready! Now ask the student to write /start to the bot.", emojiInfo))
}

// handleSetPrice sets the individual lesson price for a student.
// Usage: /set_price StudentID Price
func (b *Bot) handleSetPrice(message *tgbotapi.Message) {
	args := strings.Fields(message.CommandArguments())
	if len(args) != 2 {
		b.send(message.Chat.ID, fmt.Sprintf("%s Usage: /set_price StudentID Price\nExample: /set_price 1 1500", emojiWarning))
		return
	}

	studentID, err := strconv.Atoi(args[0])
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s StudentID must be a number", emojiError))
		return
	}

	price, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Price must be a number", emojiError))
		return
	}

	if err := b.service.SetLessonPrice(studentID, price); err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error setting price: %s", emojiError, err.Error()))
		return
	}

	b.send(message.Chat.ID, fmt.Sprintf("%s Price set to %.2f for student ID %d", emojiSuccess, price, studentID))
}

// handleAddPackage adds a paid lesson package for a student.
// Usage: /add_package StudentID TotalLessons Price
func (b *Bot) handleAddPackage(message *tgbotapi.Message) {
	args := strings.Fields(message.CommandArguments())
	if len(args) != 3 {
		b.send(message.Chat.ID, fmt.Sprintf("%s Usage: /add_package StudentID TotalLessons Price\nExample: /add_package 1 8 10000", emojiWarning))
		return
	}

	studentID, err := strconv.Atoi(args[0])
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s StudentID must be a number", emojiError))
		return
	}

	totalLessons, err := strconv.Atoi(args[1])
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s TotalLessons must be a number", emojiError))
		return
	}

	price, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Price must be a number", emojiError))
		return
	}

	id, err := b.service.AddLessonPackage(studentID, totalLessons, price)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error adding package: %s", emojiError, err.Error()))
		return
	}

	b.send(message.Chat.ID, fmt.Sprintf("%s Package of %d lessons added with ID %d", emojiSuccess, totalLessons, id))
}

// handleLessonDone marks a lesson as conducted for a student.
// If only 1 lesson remains after this — sends a payment reminder.
// Usage: /lesson_done StudentID
func (b *Bot) handleLessonDone(message *tgbotapi.Message) {
	keyboard, err := b.buildStudentKeyboard("lesson_done")
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error getting students: %s", emojiError, err.Error()))
		return
	}

	b.sendWithKeyboard(message.Chat.ID, "Select a student:", keyboard)
}

// handleBalance shows the active package balance for a student.
// Usage: /balance StudentID
func (b *Bot) handleBalance(message *tgbotapi.Message) {
	keyboard, err := b.buildStudentKeyboard("balance")
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error getting students: %s", emojiError, err.Error()))
		return
	}

	b.sendWithKeyboard(message.Chat.ID, "Select a student:", keyboard)
}

// handleMyHomeworks shows homework tasks for the student who sent the command.
func (b *Bot) handleMyHomeworks(message *tgbotapi.Message) {
	student, err := b.service.GetStudentByTelegramID(message.From.ID)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s You are not linked to any student. Ask your teacher to link your account.", emojiError))
		return
	}

	homeworks, err := b.service.GetHomeworksByStudent(student.ID)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s Error getting homeworks: %s", emojiError, err.Error()))
		return
	}

	if len(homeworks) == 0 {
		b.send(message.Chat.ID, fmt.Sprintf("%s No homework yet!", emojiInfo))
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s Your homeworks:\n\n", emojiBook))
	for _, hw := range homeworks {
		status := emojiCross
		if hw.Done {
			status = emojiCheck
		}
		sb.WriteString(fmt.Sprintf("%s %s (added: %s)\n", status, hw.Task, hw.CreatedAt))
	}

	b.send(message.Chat.ID, sb.String())
}

// handleMyLessons shows the remaining lessons in the active package for the student.
func (b *Bot) handleMyLessons(message *tgbotapi.Message) {
	student, err := b.service.GetStudentByTelegramID(message.From.ID)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s You are not linked to any student. Ask your teacher to link your account.", emojiError))
		return
	}

	pkg, err := b.service.GetActivePackage(student.ID)
	if err != nil {
		b.send(message.Chat.ID, fmt.Sprintf("%s No active lesson package. Contact your teacher.", emojiInfo))
		return
	}

	remaining := pkg.TotalLessons - pkg.UsedLessons
	b.send(message.Chat.ID, fmt.Sprintf(
		"%s Remaining lessons: %d out of %d",
		emojiCalendar, remaining, pkg.TotalLessons,
	))
}

// handleCallback processes inline keyboard button clicks.
// Callback data format: "action:studentID", e.g. "lesson_done:2"
func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) {
	// always answer the callback to remove the loading spinner in Telegram
	b.api.Request(tgbotapi.NewCallback(callback.ID, ""))

	parts := strings.SplitN(callback.Data, ":", 2)
	if len(parts) != 2 {
		return
	}

	action := parts[0]
	studentID, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	switch action {
	case "lesson_done":
		b.callbackLessonDone(callback, studentID)
	case "balance":
		b.callbackBalance(callback, studentID)
	case "homeworks":
		b.callbackHomeworks(callback, studentID)
	}
}

func (b *Bot) callbackLessonDone(callback *tgbotapi.CallbackQuery, studentID int) {
	remind, err := b.service.ConductLesson(studentID)
	if err != nil {
		b.send(callback.Message.Chat.ID, fmt.Sprintf("%s Error: %s", emojiError, err.Error()))
		return
	}

	b.send(callback.Message.Chat.ID, fmt.Sprintf("%s Lesson marked as conducted for student ID %d", emojiSuccess, studentID))

	if remind {
		// only 1 lesson left — remind teacher about payment
		b.send(callback.Message.Chat.ID, fmt.Sprintf("%s Reminder: only 1 lesson left in the package for student ID %d. Time to discuss the next payment!", emojiWarning, studentID))
	}
}

func (b *Bot) callbackBalance(callback *tgbotapi.CallbackQuery, studentID int) {
	pkg, err := b.service.GetActivePackage(studentID)
	if err != nil {
		b.send(callback.Message.Chat.ID, fmt.Sprintf("%s No active package for student ID %d", emojiInfo, studentID))
		return
	}

	remaining := pkg.TotalLessons - pkg.UsedLessons
	b.send(callback.Message.Chat.ID, fmt.Sprintf(
		"%s Balance for student ID %d:\nPackage: %d lessons\nUsed: %d\nRemaining: %d\nPrice paid: %.2f",
		emojiInfo, studentID, pkg.TotalLessons, pkg.UsedLessons, remaining, pkg.Price,
	))
}

func (b *Bot) callbackHomeworks(callback *tgbotapi.CallbackQuery, studentID int) {
	homeworks, err := b.service.GetHomeworksByStudent(studentID)
	if err != nil {
		b.send(callback.Message.Chat.ID, fmt.Sprintf("%s Error getting homeworks: %s", emojiError, err.Error()))
		return
	}

	if len(homeworks) == 0 {
		b.send(callback.Message.Chat.ID, fmt.Sprintf("%s No homeworks for this student yet.", emojiInfo))
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s Homeworks:\n\n", emojiBook))
	for _, hw := range homeworks {
		status := emojiCross
		if hw.Done {
			status = emojiCheck
		}
		sb.WriteString(fmt.Sprintf("ID %d - %s - %s (added: %s)\n", hw.ID, status, hw.Task, hw.CreatedAt))
	}

	b.send(callback.Message.Chat.ID, sb.String())
}
