package router

import (
	"krypton/pkg/auth"
	"krypton/wallet/internal/application/services"
	"krypton/wallet/internal/delivery/http/handlers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Setup(e *gin.Engine, walletService services.WalletService, tokenManager *auth.TokenManager, logger *zap.Logger) {
	walletHandler := handlers.NewWalletHandler(walletService, logger)

	v1 := e.Group("/v1")
	v1.Use(tokenManager.AuthMiddleware())
	{
		// Accounts
		v1.GET("/accounts", walletHandler.GetAccounts)
		v1.GET("/accounts/:asset_symbol", walletHandler.GetAccountByAsset)

		// v1.POST("/addresses", walletHandler.GetOrCreateDepositAddress)
		// v1.POST("/withdrawals", walletHandler.RequestWithdrawal)
	}
}
