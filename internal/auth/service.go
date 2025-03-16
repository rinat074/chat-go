package auth

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"encoding/base64"

	"chat-app/internal/db"
	"chat-app/internal/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("неверные учетные данные")
	ErrUserExists         = errors.New("пользователь уже существует")
	jwtSecret             = []byte("ваш-секретный-ключ") // В продакшене должен быть безопасно сохранен
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour // 7 дней
)

type Service struct {
	db *db.Database
}

func NewService(db *db.Database) *Service {
	return &Service{db: db}
}

func (s *Service) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	// Проверка существования пользователя
	var exists bool
	err := s.db.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)",
		req.Username, req.Email).Scan(&exists)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrUserExists
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создание пользователя
	var user models.User
	err = s.db.Pool.QueryRow(ctx,
		`INSERT INTO users (username, email, password, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, username, email, created_at, updated_at`,
		req.Username, req.Email, string(hashedPassword), time.Now(), time.Now(),
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Создание JWT токена
	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) Login(ctx context.Context, req models.LoginRequest, userAgent, ip string) (*models.AuthResponse, error) {
	var user models.User
	var hashedPassword string

	err := s.db.Pool.QueryRow(ctx,
		`SELECT id, username, email, password, created_at, updated_at 
		FROM users WHERE username = $1`,
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Создание токенов
	tokenPair, err := s.generateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// Сохранение сессии refresh токена
	refreshSession := models.RefreshSession{
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		UserAgent:    userAgent,
		IP:           ip,
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
		CreatedAt:    time.Now(),
	}

	_, err = s.db.Pool.Exec(ctx,
		`INSERT INTO refresh_sessions (user_id, refresh_token, user_agent, ip, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		refreshSession.UserID, refreshSession.RefreshToken, refreshSession.UserAgent,
		refreshSession.IP, refreshSession.ExpiresAt, refreshSession.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Tokens: tokenPair,
		User:   user,
	}, nil
}

func (s *Service) generateToken(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *Service) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неожиданный метод подписи")
		}
		return jwtSecret, nil
	})
}

// Генерация пары токенов
func (s *Service) generateTokenPair(userID int64, username string) (*models.TokenPair, error) {
	// Генерируем access token
	accessClaims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(accessTokenTTL).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// Генерируем refresh token
	refreshToken := make([]byte, 32)
	if _, err := rand.Read(refreshToken); err != nil {
		return nil, err
	}
	refreshTokenString := base64.URLEncoding.EncodeToString(refreshToken)

	// Время истечения токена
	expiresAt := time.Now().Add(accessTokenTTL)

	return &models.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshToken обновляет пару токенов
func (s *Service) RefreshToken(ctx context.Context, refreshToken, userAgent, ip string) (*models.TokenPair, error) {
	// Проверяем существование сессии с таким refresh токеном
	var session models.RefreshSession
	var username string

	err := s.db.Pool.QueryRow(ctx,
		`SELECT rs.id, rs.user_id, rs.refresh_token, rs.user_agent, rs.ip, rs.expires_at, rs.created_at, u.username
		FROM refresh_sessions rs
		JOIN users u ON rs.user_id = u.id
		WHERE rs.refresh_token = $1 AND rs.expires_at > NOW()`,
		refreshToken,
	).Scan(&session.ID, &session.UserID, &session.RefreshToken, &session.UserAgent,
		&session.IP, &session.ExpiresAt, &session.CreatedAt, &username)

	if err != nil {
		return nil, errors.New("недействительный refresh токен")
	}

	// Удаляем старую сессию
	_, err = s.db.Pool.Exec(ctx, "DELETE FROM refresh_sessions WHERE id = $1", session.ID)
	if err != nil {
		return nil, err
	}

	// Генерируем новую пару токенов
	tokenPair, err := s.generateTokenPair(session.UserID, username)
	if err != nil {
		return nil, err
	}

	// Сохраняем новую сессию
	newSession := models.RefreshSession{
		UserID:       session.UserID,
		RefreshToken: tokenPair.RefreshToken,
		UserAgent:    userAgent,
		IP:           ip,
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
		CreatedAt:    time.Now(),
	}

	_, err = s.db.Pool.Exec(ctx,
		`INSERT INTO refresh_sessions (user_id, refresh_token, user_agent, ip, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		newSession.UserID, newSession.RefreshToken, newSession.UserAgent,
		newSession.IP, newSession.ExpiresAt, newSession.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

// Logout - выйти из системы, удалив refresh токен
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	_, err := s.db.Pool.Exec(ctx, "DELETE FROM refresh_sessions WHERE refresh_token = $1", refreshToken)
	return err
}
