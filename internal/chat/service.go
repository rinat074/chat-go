package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"chat-app/internal/cache"
	"chat-app/internal/db"
	"chat-app/internal/models"

	"github.com/jackc/pgx/v5"
)

var (
	ErrGroupNotFound  = errors.New("группа не найдена")
	ErrNotGroupMember = errors.New("пользователь не является членом группы")
)

type Service struct {
	db    *db.Database
	cache *cache.Cache
}

func NewService(db *db.Database, cache *cache.Cache) *Service {
	return &Service{db: db, cache: cache}
}

// SaveMessage сохраняет сообщение в базе данных
func (s *Service) SaveMessage(ctx context.Context, msg models.Message) (*models.Message, error) {
	var savedMsg models.Message

	query := `INSERT INTO messages (type, content, user_id, username, receiver_id, group_id, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) 
	          RETURNING id, type, content, user_id, username, receiver_id, group_id, created_at`

	err := s.db.Pool.QueryRow(ctx, query,
		msg.Type, msg.Content, msg.UserID, msg.Username, msg.ReceiverID, msg.GroupID, msg.CreatedAt,
	).Scan(&savedMsg.ID, &savedMsg.Type, &savedMsg.Content, &savedMsg.UserID,
		&savedMsg.Username, &savedMsg.ReceiverID, &savedMsg.GroupID, &savedMsg.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Инвалидируем кэш при добавлении новых сообщений
	switch msg.Type {
	case models.PublicMessage:
		s.cache.DeleteByPattern(ctx, "messages:public:*")
	case models.PrivateMessage:
		if msg.ReceiverID != nil {
			// Удаляем кэш для обоих участников личной переписки
			s.cache.DeleteByPattern(ctx, fmt.Sprintf("messages:private:%d:%d:*", msg.UserID, *msg.ReceiverID))
			s.cache.DeleteByPattern(ctx, fmt.Sprintf("messages:private:%d:%d:*", *msg.ReceiverID, msg.UserID))
		}
	case models.GroupMessage:
		if msg.GroupID != nil {
			s.cache.DeleteByPattern(ctx, fmt.Sprintf("messages:group:%d:*", *msg.GroupID))
		}
	}

	return &savedMsg, nil
}

// GetPublicMessages возвращает историю публичных сообщений
func (s *Service) GetPublicMessages(ctx context.Context, limit, offset int) ([]models.Message, error) {
	// Пытаемся получить данные из кэша
	cacheKey := fmt.Sprintf("messages:public:%d:%d", limit, offset)
	messages, err := s.cache.GetMessages(ctx, cacheKey)
	if err == nil {
		return messages, nil
	}

	// Если данных в кэше нет, получаем из БД
	rows, err := s.db.Pool.Query(ctx,
		`SELECT id, type, content, user_id, username, receiver_id, group_id, created_at 
		FROM messages 
		WHERE type = 'public'
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages, err = s.scanMessages(rows)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш на 5 минут
	s.cache.SetMessages(ctx, cacheKey, messages, 5*time.Minute)

	return messages, nil
}

// GetPrivateMessages возвращает личные сообщения между двумя пользователями
func (s *Service) GetPrivateMessages(ctx context.Context, userID, otherUserID int64, limit, offset int) ([]models.Message, error) {
	rows, err := s.db.Pool.Query(ctx,
		`SELECT id, type, content, user_id, username, receiver_id, group_id, created_at 
		FROM messages 
		WHERE type = 'private' AND (
			(user_id = $1 AND receiver_id = $2) OR 
			(user_id = $2 AND receiver_id = $1)
		)
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4`,
		userID, otherUserID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanMessages(rows)
}

// GetGroupMessages возвращает сообщения группы
func (s *Service) GetGroupMessages(ctx context.Context, groupID, userID int64, limit, offset int) ([]models.Message, error) {
	// Проверяем, является ли пользователь членом группы
	var isMember bool
	err := s.db.Pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)`,
		groupID, userID,
	).Scan(&isMember)

	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, ErrNotGroupMember
	}

	rows, err := s.db.Pool.Query(ctx,
		`SELECT id, type, content, user_id, username, receiver_id, group_id, created_at 
		FROM messages 
		WHERE type = 'group' AND group_id = $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`,
		groupID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanMessages(rows)
}

// CreateGroup создает новую группу
func (s *Service) CreateGroup(ctx context.Context, group models.Group) (*models.Group, error) {
	var savedGroup models.Group

	err := s.db.Pool.QueryRow(ctx,
		`INSERT INTO groups (name, description, owner_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, name, description, owner_id, created_at, updated_at`,
		group.Name, group.Description, group.OwnerID, group.CreatedAt, group.UpdatedAt,
	).Scan(&savedGroup.ID, &savedGroup.Name, &savedGroup.Description,
		&savedGroup.OwnerID, &savedGroup.CreatedAt, &savedGroup.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Добавляем создателя как члена группы (роль владельца)
	_, err = s.db.Pool.Exec(ctx,
		`INSERT INTO group_members (group_id, user_id, role, joined_at) 
		VALUES ($1, $2, 'owner', $3)`,
		savedGroup.ID, savedGroup.OwnerID, savedGroup.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &savedGroup, nil
}

// AddUserToGroup добавляет пользователя в группу
func (s *Service) AddUserToGroup(ctx context.Context, groupID, userID, adminID int64) error {
	// Проверяем, имеет ли adminID права для добавления пользователей (owner или admin)
	var adminRole string
	err := s.db.Pool.QueryRow(ctx,
		`SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2`,
		groupID, adminID,
	).Scan(&adminRole)

	if err != nil {
		return err
	}

	if adminRole != "owner" && adminRole != "admin" {
		return errors.New("недостаточно прав")
	}

	// Добавляем пользователя в группу
	_, err = s.db.Pool.Exec(ctx,
		`INSERT INTO group_members (group_id, user_id, role, joined_at) 
		VALUES ($1, $2, 'member', NOW())`,
		groupID, userID,
	)

	return err
}

// Вспомогательный метод для сканирования сообщений из результата запроса
func (s *Service) scanMessages(rows pgx.Rows) ([]models.Message, error) {
	messages := []models.Message{}
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.Type, &msg.Content, &msg.UserID,
			&msg.Username, &msg.ReceiverID, &msg.GroupID, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
