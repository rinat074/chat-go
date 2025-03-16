package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database представляет подключение к базе данных
type Database struct {
	Pool *pgxpool.Pool
}

// NewConnection создает новое подключение к базе данных
func NewConnection(connStr string) (*Database, error) {
	// Создаем пул подключений
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул подключений: %w", err)
	}

	// Проверяем соединение
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	return &Database{Pool: pool}, nil
}

// Close закрывает подключение к базе данных
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
