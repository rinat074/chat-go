package server

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/gochat/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"chat-app/services/auth-service/internal/models"
	"chat-app/services/auth-service/internal/service"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	service *service.AuthService
}

func RegisterAuthServer(s *grpc.Server, svc *service.AuthService) {
	auth.RegisterAuthServiceServer(s, &AuthServer{
		service: svc,
	})
}

// Register регистрирует нового пользователя
func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.AuthResponse, error) {
	// Преобразуем запрос в модель домена
	registerReq := models.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Вызываем сервис для регистрации
	resp, err := s.service.Register(ctx, registerReq)
	if err != nil {
		return nil, err
	}

	// Преобразуем ответ в protobuf
	protoResp := &auth.AuthResponse{
		User: &auth.User{
			Id:        resp.User.ID,
			Username:  resp.User.Username,
			Email:     resp.User.Email,
			CreatedAt: timestamppb.New(resp.User.CreatedAt),
			UpdatedAt: timestamppb.New(resp.User.UpdatedAt),
		},
	}

	// Если есть токен, добавляем его
	if resp.Token != "" {
		protoResp.Tokens = &auth.TokenPair{
			AccessToken: resp.Token,
		}
	} else if resp.Tokens != nil {
		protoResp.Tokens = &auth.TokenPair{
			AccessToken:  resp.Tokens.AccessToken,
			RefreshToken: resp.Tokens.RefreshToken,
			ExpiresAt:    timestamppb.New(resp.Tokens.ExpiresAt),
		}
	}

	return protoResp, nil
}

// Login выполняет вход пользователя
func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.AuthResponse, error) {
	// Преобразуем запрос в модель домена
	loginReq := models.LoginRequest{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: req.UserAgent,
		IP:        req.Ip,
	}

	// Вызываем сервис для входа
	resp, err := s.service.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	// Преобразуем ответ в protobuf
	protoResp := &auth.AuthResponse{
		User: &auth.User{
			Id:        resp.User.ID,
			Username:  resp.User.Username,
			Email:     resp.User.Email,
			CreatedAt: timestamppb.New(resp.User.CreatedAt),
			UpdatedAt: timestamppb.New(resp.User.UpdatedAt),
		},
	}

	// Если есть токен, добавляем его
	if resp.Token != "" {
		protoResp.Tokens = &auth.TokenPair{
			AccessToken: resp.Token,
		}
	} else if resp.Tokens != nil {
		protoResp.Tokens = &auth.TokenPair{
			AccessToken:  resp.Tokens.AccessToken,
			RefreshToken: resp.Tokens.RefreshToken,
			ExpiresAt:    timestamppb.New(resp.Tokens.ExpiresAt),
		}
	}

	return protoResp, nil
}

// RefreshToken обновляет пару токенов
func (s *AuthServer) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.TokenPair, error) {
	// Вызываем сервис для обновления токена
	tokenPair, err := s.service.RefreshToken(ctx, req.RefreshToken, req.UserAgent, req.Ip)
	if err != nil {
		return nil, err
	}

	// Преобразуем ответ в protobuf
	return &auth.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    timestamppb.New(tokenPair.ExpiresAt),
	}, nil
}

// Logout выполняет выход пользователя
func (s *AuthServer) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// Вызываем сервис для выхода
	err := s.service.Logout(ctx, req.RefreshToken)
	if err != nil {
		return &auth.LogoutResponse{Success: false}, err
	}

	return &auth.LogoutResponse{Success: true}, nil
}

// ValidateToken проверяет валидность токена
func (s *AuthServer) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	// Вызываем сервис для проверки токена
	token, err := s.service.ValidateToken(req.Token)
	if err != nil {
		return &auth.ValidateTokenResponse{Valid: false}, nil
	}

	// Проверяем, что токен валиден
	if !token.Valid {
		return &auth.ValidateTokenResponse{Valid: false}, nil
	}

	// Получаем данные из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &auth.ValidateTokenResponse{Valid: false}, nil
	}

	// Получаем ID и имя пользователя
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return &auth.ValidateTokenResponse{Valid: false}, nil
	}

	username, ok := claims["username"].(string)
	if !ok {
		return &auth.ValidateTokenResponse{Valid: false}, nil
	}

	return &auth.ValidateTokenResponse{
		Valid:    true,
		UserId:   int64(userID),
		Username: username,
	}, nil
}
