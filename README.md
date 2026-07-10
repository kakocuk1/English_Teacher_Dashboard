# English Teacher Dashboard

A Telegram bot for managing English students, homework, schedules, and lesson payments.

Built as a portfolio project to demonstrate Go development skills.

## Features

- Add and manage students with individual lesson prices
- Assign homework and track completion
- Schedule lessons by day and time
- Track lesson packages and payments
- Automatic payment reminder when 1 lesson remains
- Students can link their Telegram account to view their homework and remaining lessons

## Tech Stack

- Go 1.25.6
- SQLite (modernc.org/sqlite — pure Go, no CGO)
- Telegram Bot API (go-telegram-bot-api/v5)
- Docker

## Project Structure
teacher-dashboard/
├── cmd/
│   └── bot/          # entry point
├── internal/
│   ├── model/        # data structures
│   ├── storage/      # SQLite layer
│   ├── service/      # business logic
│   └── bot/          # Telegram handlers
└── Dockerfile

## Getting Started

### Prerequisites

- Go 1.25.6+
- Telegram bot token from [@BotFather](https://t.me/BotFather)

### Run locally

```bash
# clone the repo
git clone https://github.com/kakocuk1/English_Teacher_Dashboard.git
cd English_Teacher_Dashboard

# set your token
echo "TELEGRAM_TOKEN=your_token_here" > .env

# run
source .env && go run ./cmd/bot/
```

### Run with Docker

```bash
docker build -t teacher-dashboard:v1 .
docker run -e TELEGRAM_TOKEN=your_token_here teacher-dashboard:v1
```

## Bot Commands
### Teacher
| Command | Description |
|---|---|
| `/add_student Name Level` | Add a student (e.g. `/add_student John B2`) |
| `/students` | List all students |
| `/link_student StudentID` | Link next student who writes /start |
| `/set_price StudentID Price` | Set individual lesson price |
| `/add_package StudentID Lessons Price` | Add paid lesson package |
| `/lesson_done StudentID` | Mark lesson as conducted |
| `/balance StudentID` | Show remaining lessons |
| `/add_homework StudentID Task` | Assign homework |
| `/homeworks StudentID` | View student homework |
| `/done HomeworkID` | Mark homework as done |
| `/add_lesson StudentID Day Time` | Add lesson to schedule |
| `/schedule` | View full schedule |

### Student
| Command | Description |
|---|---|
| `/my_homeworks` | View your homework |
| `/my_lessons` | View remaining lessons |

## Architecture

Business logic lives in `internal/service/` behind a `Storage` interface.
The Telegram bot is just one transport layer — a web UI can be added later without rewriting the core logic.

## License

MIT