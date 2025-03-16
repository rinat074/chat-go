# Auth Service

Сервис аутентификации и авторизации для чат-приложения GoChat.

## Функциональность

- Регистрация новых пользователей
- Аутентификация пользователей
- Управление JWT токенами
- Обновление токенов доступа
- Выход из системы
- Проверка валидности токенов

## API

Сервис предоставляет следующие gRPC методы:

### Register

Регистрация нового пользователя.

```protobuf
rpc Register(RegisterRequest) returns (AuthResponse)
```

### Login

Вход пользователя в систему.

```protobuf
rpc Login(LoginRequest) returns (AuthResponse)
```

### RefreshToken

Обновление токена доступа.

```protobuf
rpc RefreshToken(RefreshTokenRequest) returns (TokenPair)
```

### Logout

Выход пользователя из системы.

```protobuf
rpc Logout(LogoutRequest) returns (LogoutResponse)
```

### ValidateToken

Проверка валидности токена.

```protobuf
rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse)
```

## Конфигурация

Сервис настраивается через переменные окружения:

- `DATABASE_URL` - URL подключения к PostgreSQL (по умолчанию: `postgres://postgres:postgres@postgres:5432/chatapp`)
- `GRPC_SERVER_ADDRESS` - адрес gRPC сервера (по умолчанию: `:50051`)
- `JWT_SECRET` - секретный ключ для JWT токенов

## Структура базы данных

### Таблица users

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Таблица refresh_sessions

```sql
CREATE TABLE refresh_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    user_agent VARCHAR(255),
    ip VARCHAR(45),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);
```

## Запуск

### Запуск с помощью Docker

```bash
docker-compose up auth-service
```

### Запуск без Docker

```bash
cd services/auth-service
go build -o ../../bin/auth-service ./cmd
../../bin/auth-service
``` 