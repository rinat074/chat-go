-- Создание базы данных
CREATE DATABASE chatapp;

\c chatapp;

-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Таблица сессий обновления токенов
CREATE TABLE refresh_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    user_agent VARCHAR(255),
    ip VARCHAR(45),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- Таблица сообщений
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    type VARCHAR(10) NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    group_id INTEGER,
    created_at TIMESTAMP NOT NULL
);

-- Таблица групп
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Таблица участников групп
CREATE TABLE group_members (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    UNIQUE(group_id, user_id)
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX idx_messages_group_id ON messages(group_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_group_members_group_id ON group_members(group_id);
CREATE INDEX idx_group_members_user_id ON group_members(user_id);
CREATE INDEX idx_refresh_sessions_user_id ON refresh_sessions(user_id);
CREATE INDEX idx_refresh_sessions_refresh_token ON refresh_sessions(refresh_token);

-- Создание тестовых пользователей
INSERT INTO users (username, email, password, created_at, updated_at)
VALUES 
    ('user1', 'user1@example.com', '$2a$10$1qAz2wSx3eCc1/ps3TEeP.uAGI/7KcC9gEwW0s0T7UY6Ldl8EiHAS', NOW(), NOW()), -- password: password1
    ('user2', 'user2@example.com', '$2a$10$37nMN1B9jlMx/ZNUDKNlJeVs9m.QpCCK1RYlQpGnsb3E1Q6fK5LK2', NOW(), NOW()); -- password: password2

-- Создание тестовой группы
INSERT INTO groups (name, description, owner_id, created_at, updated_at)
VALUES ('Общий чат', 'Чат для всех пользователей', 1, NOW(), NOW());

-- Добавление пользователей в группу
INSERT INTO group_members (group_id, user_id, created_at)
VALUES 
    (1, 1, NOW()),
    (1, 2, NOW());

-- Создание тестовых сообщений
INSERT INTO messages (type, content, user_id, receiver_id, group_id, created_at)
VALUES 
    ('public', 'Привет всем!', 1, NULL, NULL, NOW()),
    ('private', 'Привет, как дела?', 1, 2, NULL, NOW()),
    ('group', 'Всем привет в группе!', 1, NULL, 1, NOW()),
    ('group', 'И тебе привет!', 2, NULL, 1, NOW()); 