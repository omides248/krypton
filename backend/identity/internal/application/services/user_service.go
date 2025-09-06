package services

import (
	"context"
	"errors"
	"krypton/identity/internal/domain"
	"krypton/pkg/auth"

	"go.uber.org/zap"
)

type userService struct {
	userRepo     domain.UserRepository
	tokenManager *auth.TokenManager
	logger       *zap.Logger
}

func NewUserService(userRepo domain.UserRepository, tokenManager *auth.TokenManager, logger *zap.Logger) UserService {
	return &userService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		logger:       logger.Named("user_service"),
	}
}

func (s *userService) Register(ctx context.Context, email, plainPassword string) (*domain.User, error) {
	s.logger.Info("processing user registration", zap.String("email", email))

	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		s.logger.Warn("registration failed: email already exists", zap.String("email", email))
		return nil, domain.ErrEmailAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		s.logger.Error("failed to check for existing email", zap.Error(err))
		return nil, err
	}

	newUser, err := domain.NewUser(email, plainPassword)
	if err != nil {
		s.logger.Error("failed to create new user object", zap.Error(err))
		return nil, err
	}

	if err := s.userRepo.Save(ctx, newUser); err != nil {
		s.logger.Error("failed to save new user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("user registered successfully", zap.String("user_id", string(newUser.ID)))

	// TODO propagate domain event

	return newUser, nil
}

func (s *userService) Login(ctx context.Context, email, plainPassword string) (user *domain.User, accessToken string, refreshToken string, err error) {
	s.logger.Info("processing user login", zap.String("email", email))

	user, err = s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("login failed: cannot find user", zap.String("email", email), zap.Error(err))
		return nil, "", "", domain.ErrInvalidPassword
	}

	if err := user.CheckPassword(plainPassword); err != nil {
		s.logger.Warn("login failed: invalid password", zap.String("email", email))
		return nil, "", "", domain.ErrInvalidPassword
	}

	accessToken, err = s.tokenManager.GenerateAccessToken(string(user.ID))
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err = s.tokenManager.GenerateRefreshToken(string(user.ID))
	if err != nil {
		return nil, "", "", err
	}

	s.logger.Info("user logged in successfully", zap.String("user_id", string(user.ID)))

	return user, accessToken, refreshToken, nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshTokenString string) (string, error) {
	s.logger.Info("processing token refresh")

	claims, err := s.tokenManager.Validate(refreshTokenString)
	if err != nil || claims.Subject != "refresh_token" {
		s.logger.Warn("refresh failed: invalid refresh token")
		return "", domain.ErrInvalidToken
	}

	_, err = s.userRepo.FindByID(ctx, domain.UserID(claims.UserID))
	if err != nil {
		s.logger.Warn("refresh failed: user not found", zap.String("user_id", claims.UserID))
		return "", domain.ErrUserNotFound
	}

	newAccessToken, err := s.tokenManager.GenerateAccessToken(claims.UserID)
	if err != nil {
		s.logger.Error("refresh failed: could not generate new access token", zap.Error(err))
		return "", err
	}

	s.logger.Info("token refreshed successfully", zap.String("user_id", claims.UserID))
	return newAccessToken, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	s.logger.Info("fetching user profile", zap.String("user_id", string(userID)))

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to get user by id", zap.String("user_id", string(userID)), zap.Error(err))
		return nil, err
	}

	return user, nil
}
