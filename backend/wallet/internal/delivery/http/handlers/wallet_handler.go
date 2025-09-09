package handlers

import (
	"krypton/pkg/contextkeys"
	"krypton/wallet/internal/application/services"
	"krypton/wallet/internal/delivery/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WalletHandler struct {
	service services.WalletService
	logger  *zap.Logger
}

func NewWalletHandler(service services.WalletService, logger *zap.Logger) *WalletHandler {
	return &WalletHandler{
		service: service,
		logger:  logger.Named("wallet_http_handler"),
	}
}

func (h *WalletHandler) GetAccounts(c *gin.Context) {
	userID, err := contextkeys.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	userAccounts, err := h.service.GetAccounts(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToAccountsResponse(userAccounts))
	return
}

func (h *WalletHandler) GetAccountByAsset(c *gin.Context) {
	userID, err := contextkeys.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	assetSymbol := c.Param("asset_symbol")

	userAccount, err := h.service.GetAccountByAsset(c.Request.Context(), userID, assetSymbol)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToAccountResponse(*userAccount))
	return
}
