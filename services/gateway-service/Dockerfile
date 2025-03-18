FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копирование go.mod и go.sum файлов
COPY go.mod go.sum* ./
COPY services/gateway-service/go.mod services/gateway-service/go.sum* ./services/gateway-service/
COPY proto/auth/go.mod proto/auth/go.sum* ./proto/auth/
COPY proto/chat/go.mod proto/chat/go.sum* ./proto/chat/

# Скачивание зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN cd services/gateway-service && \
    GOWORK=off go build -o /app/bin/gateway-service ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Установка зависимостей
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копирование бинарного файла
COPY --from=builder /app/bin/gateway-service /app/gateway-service

# Установка порта
EXPOSE 8080

# Запуск приложения
ENTRYPOINT ["/app/gateway-service"] 