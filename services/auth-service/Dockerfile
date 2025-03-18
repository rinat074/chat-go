FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копирование go.mod и go.sum файлов
COPY go.mod go.sum* ./
COPY services/auth-service/go.mod services/auth-service/go.sum* ./services/auth-service/
COPY proto/auth/go.mod proto/auth/go.sum* ./proto/auth/

# Скачивание зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN cd services/auth-service && \
    go build -o /app/bin/auth-service ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Установка зависимостей
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копирование бинарного файла
COPY --from=builder /app/bin/auth-service /app/auth-service

# Запуск приложения
ENTRYPOINT ["/app/auth-service"] 