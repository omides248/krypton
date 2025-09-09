package services

import (
	"context"
	"krypton/wallet/internal/domain"
	"math/big"

	"go.uber.org/zap"
)

type walletService struct {
	accountRepo domain.AccountRepository
	assetRepo   domain.AssetRepository
	logger      *zap.Logger
}

func NewWalletService(
	accountRepo domain.AccountRepository,
	assetRepo domain.AssetRepository,
	logger *zap.Logger) WalletService {

	return &walletService{
		accountRepo: accountRepo,
		assetRepo:   assetRepo,
		logger:      logger.Named("wallet_service"),
	}
}

func (s *walletService) GetAccounts(ctx context.Context, userID string) ([]domain.UserAccount, error) {
	s.logger.Info("Fetching all accounts for user", zap.String("userID", userID))
	return s.accountRepo.FindUserAccounts(ctx, userID)
}

func (s *walletService) GetAccountByAsset(ctx context.Context, userID string, assetSymbol string) (*domain.UserAccount, error) {
	s.logger.Info("Fetching account for user by asset",
		zap.String("userID", userID),
		zap.String("assetSymbol", assetSymbol),
	)
	panic("implement me")
}

func (s *walletService) RequestWithdrawal(ctx context.Context, userID string, assetSymbol string, amount *big.Int, toAddress string) (*domain.Withdrawal, error) {
	//TODO implement me
	panic("implement me")
}
