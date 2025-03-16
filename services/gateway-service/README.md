# Gateway Service

Gateway Service - это API-шлюз для приложения GoChat, который обрабатывает HTTP-запросы и WebSocket-соединения, взаимодействуя с Auth Service и Chat Service через gRPC.

## Функциональность

- Обработка HTTP-запросов для аутентификации и управления чатами
- Поддержка WebSocket-соединений для обмена сообщениями в реальном времени
- Маршрутизация запросов к соответствующим микросервисам
- Проверка JWT-токенов для защиты API-эндпоинтов
- Управление сессиями пользователей

## API Endpoints

### Аутентификация

- `POST /api/auth/register` - Регистрация нового пользователя
  - Запрос: `{ "username": "string", "password": "string", "email": "string" }`
  - Ответ: `{ "user_id": "string", "username": "string", "access_token": "string", "refresh_token": "string" }`

- `POST /api/auth/login` - Вход пользователя
  - Запрос: `{ "username": "string", "password": "string" }`
  - Ответ: `{ "user_id": "string", "username": "string", "access_token": "string", "refresh_token": "string" }`

- `POST /api/auth/refresh` - Обновление токена доступа
  - Запрос: `{ "refresh_token": "string" }`
  - Ответ: `{ "access_token": "string", "refresh_token": "string" }`

- `POST /api/auth/logout` - Выход пользователя
  - Запрос: `{ "refresh_token": "string" }`
  - Ответ: `{ "success": true }`

### Чаты

- `GET /api/chats/public` - Получение публичных сообщений
  - Ответ: `{ "messages": [{ "id": "string", "sender_id": "string", "sender_name": "string", "content": "string", "timestamp": "string" }] }`

- `GET /api/chats/private/:user_id` - Получение приватных сообщений с пользователем
  - Ответ: `{ "messages": [{ "id": "string", "sender_id": "string", "sender_name": "string", "content": "string", "timestamp": "string" }] }`

- `GET /api/chats/group/:group_id` - Получение сообщений группы
  - Ответ: `{ "messages": [{ "id": "string", "sender_id": "string", "sender_name": "string", "content": "string", "timestamp": "string" }] }`

- `POST /api/chats/group` - Создание новой группы
  - Запрос: `{ "name": "string", "user_ids": ["string"] }`
  - Ответ: `{ "group_id": "string", "name": "string" }`

- `POST /api/chats/group/:group_id/user` - Добавление пользователя в группу
  - Запрос: `{ "user_id": "string" }`
  - Ответ: `{ "success": true }`

### WebSocket

- `GET /ws` - WebSocket-соединение для обмена сообщениями в реальном времени
  - Параметры запроса: `?token=JWT_TOKEN`
  - Формат сообщений:
    ```json
    {
      "type": "message",
      "data": {
        "chat_type": "public|private|group",
        "recipient_id": "string", // для private и group чатов
        "content": "string"
      }
    }
    ```

## Конфигурация

Сервис можно настроить с помощью следующих переменных окружения:

- `HTTP_SERVER_ADDRESS` - Адрес HTTP-сервера (по умолчанию `:8080`)
- `AUTH_SERVICE_ADDRESS` - Адрес Auth Service (по умолчанию `localhost:50051`)
- `CHAT_SERVICE_ADDRESS` - Адрес Chat Service (по умолчанию `localhost:50052`)
- `REDIS_URL` - URL для подключения к Redis (по умолчанию `localhost:6379`)

## Запуск сервиса

### Без Docker

```bash
cd services/gateway-service
go run cmd/main.go
```

### С Docker

```bash
docker-compose up gateway-service
```

или для запуска всего приложения:

```bash
docker-compose up
``` 