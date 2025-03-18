FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копирование go.mod и go.sum файлов
COPY go.mod go.sum* ./
COPY services/chat-service/go.mod services/chat-service/go.sum* ./services/chat-service/
COPY proto/chat/go.mod proto/chat/go.sum* ./proto/chat/

# Скачивание зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN cd services/chat-service && \
    go build -o /app/bin/chat-service ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Установка зависимостей
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копирование бинарного файла
COPY --from=builder /app/bin/chat-service /app/chat-service

# Запуск приложения
ENTRYPOINT ["/app/chat-service"] 