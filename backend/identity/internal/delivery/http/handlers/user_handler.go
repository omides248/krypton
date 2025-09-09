package handlers

import (
	"errors"
	"krypton/identity/internal/application/services"
	"krypton/identity/internal/delivery/http/dto"
	"krypton/identity/internal/domain"
	"krypton/pkg/minio"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	service      services.UserService
	minioService *minio.Service
	logger       *zap.Logger
}

func NewUserHandler(service services.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger.Named("user_http_handler"),
	}
}

// GetMyProfile godoc
// @Summary      Get user profile
// @Description  Retrieves the profile information of the authenticated user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  dto.ErrorResponse "Unauthorized"
// @Failure      500  {object}  dto.ErrorResponse "Internal Server Error"
// @Security     ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @Router       /users/me [get]
func (h *UserHandler) GetMyProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		_ = c.Error(errors.New("userID not found in context"))
		return
	}

	userID := domain.UserID(userIDVal.(string))

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}
