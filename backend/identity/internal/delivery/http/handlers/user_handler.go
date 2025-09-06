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
