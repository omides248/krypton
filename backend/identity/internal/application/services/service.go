package services

import (
	"context"
	"krypton/identity/internal/domain"
	"krypton/pkg/auth"

	"go.uber.org/zap"
)

type Service struct {
	UserService UserService
}

func NewService(userRepo domain.UserRepository, tokenManager *auth.TokenManager, logger *zap.Logger) *Service {
	userSvc := NewUserService(userRepo, tokenManager, logger)
	return &Service{
		UserService: userSvc,
	}
}

type UserService interface {
	Register(ctx context.Context, email, plainPassword string) (*domain.User, error)
	Login(ctx context.Context, email, plainPassword string) (user *domain.User, accessToken string, refreshToken string, err error)
	RefreshToken(ctx context.Context, refreshTokenString string) (string, error)
	GetUserByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
}
