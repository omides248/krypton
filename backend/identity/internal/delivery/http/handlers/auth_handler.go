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

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account with the provided email and password.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body   dto.RegisterRequest true "User registration request"
// @Success      201  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponse "Validation Error"
// @Failure      409  {object}  dto.ErrorResponse "Email already exists"
// @Failure      500  {object}  dto.ErrorResponse "Internal Server Error"
// @Router       /auth/register [post]
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

// Login godoc
// @Summary      Log in a user
// @Description  Authenticates a user and returns access and refresh tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body   dto.LoginRequest true "User login request"
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  dto.ErrorResponse "Validation Error"
// @Failure      401  {object}  dto.ErrorResponse "Invalid Credentials"
// @Failure      500  {object}  dto.ErrorResponse "Internal Server Error"
// @Router       /auth/login [post]
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

// Refresh godoc
// @Summary      Refresh access token
// @Description  Generates a new access token using a valid refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body   dto.RefreshRequest true "Refresh token request"
// @Success      200  {object}  dto.RefreshResponse
// @Failure      400  {object}  dto.ErrorResponse "Validation Error"
// @Failure      401  {object}  dto.ErrorResponse "Invalid Token"
// @Failure      500  {object}  dto.ErrorResponse "Internal Server Error"
// @Router       /auth/refresh [post]
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
