# Stage 1: build the binary
FROM golang:1.25.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o teacher-dashboard ./cmd/bot/

# Stage 2: minimal runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/teacher-dashboard .

# DATABASE_URL and TELEGRAM_TOKEN are passed at runtime via environment variables
CMD ["./teacher-dashboard"]