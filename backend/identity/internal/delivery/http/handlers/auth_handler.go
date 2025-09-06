package handlers

import (
	"errors"
	"krypton/identity/internal/application/services"
	"krypton/identity/internal/delivery/http/dto"
	"krypton/identity/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service services.UserService
	logger  *zap.Logger
}

func NewAuthHandler(service services.UserService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger.Named("auth_http_handler"),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		_ = c.Error(err)
		return
	}
	if err := req.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := h.service.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("cant get user")
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.UserResponse{
		ID:    string(user.ID),
		Email: user.Email,
	})
	return
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := req.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPassword) {
			_ = c.Error(err)
			return
		}
		_ = c.Error(err)
		return
	}

	res := dto.LoginResponse{
		User: dto.UserResponse{
			ID:    string(user.ID),
			Email: user.Email,
		},
		Token: dto.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	c.JSON(http.StatusOK, res)
	return
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := req.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	newAccessToken, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.RefreshResponse{AccessToken: newAccessToken})
}
