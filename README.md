# GoChat - Микросервисное приложение для чата

[![CI/CD](https://github.com/rinat074/chat-go/actions/workflows/ci.yml/badge.svg)](https://github.com/rinat074/chat-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rinat074/chat-go)](https://goreportcard.com/report/github.com/rinat074/chat-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/rinat074/chat-go)](https://go.dev/doc/devel/release)
[![License](https://img.shields.io/github/license/rinat074/chat-go)](https://github.com/rinat074/chat-go/blob/main/LICENSE)

Это микросервисное приложение для чата, разработанное на Go. Проект демонстрирует использование микросервисной архитектуры, gRPC для межсервисного взаимодействия и RESTful API для внешних клиентов.

## Структура проекта

```
.
├── .github/workflows/     # CI/CD конфигурация для GitHub Actions
├── db/                    # Файлы для работы с базой данных
│   ├── migrations/        # Миграции SQL
│   └── init.sql           # Начальная схема базы данных
├── docker-compose.yml     # Конфигурация Docker Compose
├── Makefile               # Команды для сборки и запуска проекта
├── monitoring/            # Конфигурация для мониторинга
│   └── prometheus/        # Конфигурация Prometheus
├── pkg/                   # Общие пакеты для всех сервисов
│   ├── logger/            # Структурированное логирование
│   ├── metrics/           # Метрики Prometheus
│   ├── middleware/        # Общие middleware
│   └── validator/         # Валидация входных данных
├── proto/                 # Определения Protocol Buffers
│   ├── auth/              # Протоколы для сервиса аутентификации
│   └── chat/              # Протоколы для чат-сервиса
└── services/              # Микросервисы
    ├── auth-service/      # Сервис аутентификации
    ├── chat-service/      # Чат-сервис
    └── gateway-service/   # API-шлюз
```

## Микросервисы

### Auth Service
Отвечает за аутентификацию и авторизацию пользователей, управление JWT-токенами.

### Chat Service
Управляет чатами, сообщениями и уведомлениями.

### Gateway Service
API-шлюз, который маршрутизирует запросы к соответствующим микросервисам и предоставляет единую точку входа для клиентов.

## Технологии

- **Go**: Основной язык программирования
- **gRPC**: Для межсервисного взаимодействия
- **Protocol Buffers**: Для определения API и обмена данными
- **PostgreSQL**: Основная база данных
- **Redis**: Для кэширования и управления сессиями
- **Docker & Docker Compose**: Для контейнеризации и оркестрации
- **JWT**: Для аутентификации
- **Swagger/OpenAPI**: Для документации API
- **Prometheus & Grafana**: Для мониторинга и визуализации
- **Zap**: Для структурированного логирования
- **GitHub Actions**: Для CI/CD
- **Golang-migrate**: Для управления миграциями базы данных

## Запуск проекта

### Предварительные требования

- Go 1.21+
- Docker и Docker Compose
- Make

### Сборка и запуск

1. Клонировать репозиторий:
   ```
   git clone https://github.com/rinat074/chat-go.git
   cd chat-go
   ```

2. Запустить все сервисы с помощью Docker Compose:
   ```
   docker-compose up -d
   ```

3. Или использовать Makefile:
   ```
   make run
   ```

## Разработка

### Генерация Protocol Buffers

```
make proto
```

### Сборка отдельных сервисов

```
make build-auth
make build-chat
make build-gateway
```

### Тестирование

```
make test
make test-race  # с обнаружением гонок
```

### Линтинг

```
make lint
```

### Миграции

```
make migrate-up    # Применить миграции
make migrate-down  # Откатить миграции
```

### Сборка Docker образов

```
make docker
```

### Генерация Swagger документации

```
make swag
```

## API Endpoints

### Auth Service
- `POST /api/auth/register` - Регистрация нового пользователя
- `POST /api/auth/login` - Вход в систему
- `POST /api/auth/refresh` - Обновление токена
- `POST /api/auth/logout` - Выход из системы

### Chat Service
- `GET /api/chats` - Получение списка чатов
- `POST /api/chats` - Создание нового чата
- `GET /api/chats/{id}` - Получение информации о чате
- `GET /api/chats/{id}/messages` - Получение сообщений чата
- `POST /api/chats/{id}/messages` - Отправка сообщения в чат
- `WS /api/ws` - WebSocket соединение для реального времени

## Мониторинг

- **Prometheus**: Доступен по адресу http://localhost:9090
- **Grafana**: Доступен по адресу http://localhost:3000

## CI/CD

Проект настроен на использование GitHub Actions для автоматизации сборки, тестирования и деплоя:
- **CI**: Статический анализ кода, тестирование, проверка покрытия
- **CD**: Сборка и публикация Docker образов, автоматический деплой

## Лицензия

MIT 