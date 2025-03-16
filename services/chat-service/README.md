# Chat Service

Сервис обработки сообщений чата для приложения GoChat.

## Функциональность

- Обработка сообщений чата
- Хранение истории сообщений
- Поддержка публичных, личных и групповых чатов
- Управление группами
- Обмен сообщениями в реальном времени через WebSocket

## API

Сервис предоставляет следующие gRPC методы:

### SaveMessage

Сохранение сообщения.

```protobuf
rpc SaveMessage(Message) returns (Message)
```

### GetPublicMessages

Получение публичных сообщений.

```protobuf
rpc GetPublicMessages(GetMessagesRequest) returns (MessagesResponse)
```

### GetPrivateMessages

Получение личных сообщений между двумя пользователями.

```protobuf
rpc GetPrivateMessages(GetPrivateMessagesRequest) returns (MessagesResponse)
```

### GetGroupMessages

Получение сообщений группы.

```protobuf
rpc GetGroupMessages(GetGroupMessagesRequest) returns (MessagesResponse)
```

### CreateGroup

Создание новой группы.

```protobuf
rpc CreateGroup(CreateGroupRequest) returns (Group)
```

### AddUserToGroup

Добавление пользователя в группу.

```protobuf
rpc AddUserToGroup(AddUserToGroupRequest) returns (AddUserToGroupResponse)
```

## Конфигурация

Сервис настраивается через переменные окружения:

- `DATABASE_URL` - URL подключения к PostgreSQL (по умолчанию: `postgres://postgres:postgres@postgres:5432/chatapp`)
- `REDIS_URL` - URL подключения к Redis (по умолчанию: `redis:6379`)
- `GRPC_SERVER_ADDRESS` - адрес gRPC сервера (по умолчанию: `:50052`)

## Структура базы данных

### Таблица messages

```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    type VARCHAR(10) NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    receiver_id INTEGER,
    group_id INTEGER,
    created_at TIMESTAMP NOT NULL
);
```

### Таблица groups

```sql
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    owner_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Таблица group_members

```sql
CREATE TABLE group_members (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id),
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    UNIQUE(group_id, user_id)
);
```

## WebSocket

Сервис использует WebSocket для обмена сообщениями в реальном времени. Клиенты могут подключиться к WebSocket серверу для получения сообщений в реальном времени.

## Запуск

### Запуск с помощью Docker

```bash
docker-compose up chat-service
```

### Запуск без Docker

```bash
cd services/chat-service
go build -o ../../bin/chat-service ./cmd
../../bin/chat-service
``` 