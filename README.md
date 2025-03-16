# GoChat - Микросервисное чат-приложение на Go

GoChat - это современное чат-приложение, построенное на микросервисной архитектуре с использованием Go, gRPC, WebSocket, PostgreSQL и Redis.

## Архитектура

Приложение состоит из следующих микросервисов:

- **Auth Service** - сервис аутентификации и авторизации
- **Chat Service** - сервис обработки сообщений чата
- **Gateway Service** - API-шлюз для клиентских приложений

## Технологии

- **Go** - основной язык программирования
- **gRPC** - для коммуникации между сервисами
- **WebSocket** - для обмена сообщениями в реальном времени
- **PostgreSQL** - для хранения данных
- **Redis** - для кэширования и управления сессиями
- **Docker** - для контейнеризации
- **JWT** - для аутентификации

## Функциональность

- Регистрация и аутентификация пользователей
- Обмен сообщениями в реальном времени
- Поддержка публичных, личных и групповых чатов
- Хранение истории сообщений
- Управление пользователями и группами

## Установка и запуск

### Предварительные требования

- Go 1.21+
- Docker и Docker Compose
- PostgreSQL
- Redis

### Запуск с помощью Docker

```bash
# Клонирование репозитория
git clone https://github.com/rinat074/chat-go.git
cd chat-go

# Запуск с помощью Docker Compose
docker-compose up -d
```

### Запуск без Docker

```bash
# Клонирование репозитория
git clone https://github.com/rinat074/chat-go.git
cd chat-go

# Сборка сервисов
make build

# Запуск сервисов
./bin/auth-service
./bin/chat-service
./bin/gateway-service
```

## Структура проекта

```
.
├── cmd/                  # Точки входа для сервисов
├── proto/                # Определения Protocol Buffers
│   ├── auth/             # Протоколы аутентификации
│   └── chat/             # Протоколы чата
├── services/             # Микросервисы
│   ├── auth-service/     # Сервис аутентификации
│   ├── chat-service/     # Сервис чата
│   └── gateway-service/  # API-шлюз
├── internal/             # Общие внутренние пакеты
├── docker-compose.yml    # Конфигурация Docker Compose
└── Makefile              # Команды для сборки и запуска
```

## API

### Auth Service API

- `Register` - регистрация нового пользователя
- `Login` - вход пользователя
- `RefreshToken` - обновление токена доступа
- `Logout` - выход пользователя
- `ValidateToken` - проверка валидности токена

### Chat Service API

- `SaveMessage` - сохранение сообщения
- `GetPublicMessages` - получение публичных сообщений
- `GetPrivateMessages` - получение личных сообщений
- `GetGroupMessages` - получение групповых сообщений
- `CreateGroup` - создание группы
- `AddUserToGroup` - добавление пользователя в группу

## Разработка

### Генерация gRPC кода

```bash
make proto
```

### Сборка

```bash
make build
```

### Тестирование

```bash
make test
```

## Лицензия

MIT 