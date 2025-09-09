package services

import (
	"context"
	"krypton/wallet/internal/domain"
	"math/big"
)

type WalletService interface {
	GetAccounts(ctx context.Context, userID string) ([]domain.UserAccount, error)
	GetAccountByAsset(ctx context.Context, userID string, assetSymbol string) (*domain.UserAccount, error)
	RequestWithdrawal(ctx context.Context, userID string, assetSymbol string, amount *big.Int, toAddress string) (*domain.Withdrawal, error)
}
