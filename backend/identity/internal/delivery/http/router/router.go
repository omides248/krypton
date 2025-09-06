package router

import (
	"krypton/identity/config"
	"krypton/identity/internal/application/services"
	"krypton/identity/internal/delivery/http/handlers"
	"krypton/pkg/auth"
	"krypton/pkg/minio"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Setup(e *gin.Engine, userService services.UserService, minioService *minio.Service, cfg *config.Config, tokenManager *auth.TokenManager, logger *zap.Logger) {

	authHandler := handlers.NewAuthHandler(userService, logger)
	userHandler := handlers.NewUserHandler(userService, logger)

	v1 := e.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)

		}

		users := v1.Group("/users")
		users.Use(tokenManager.AuthMiddleware())
		{
			users.GET("/me", userHandler.GetMyProfile)
		}
	}
}
