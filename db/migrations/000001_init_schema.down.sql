-- Удаление триггеров
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_groups_updated_at ON groups;

-- Удаление функции
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаление индексов
DROP INDEX IF EXISTS idx_messages_user_id;
DROP INDEX IF EXISTS idx_messages_receiver_id;
DROP INDEX IF EXISTS idx_messages_group_id;
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_group_members_group_id;
DROP INDEX IF EXISTS idx_group_members_user_id;

-- Удаление таблиц
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users; 