package domain

import (
	"math/big"
	"time"
)

type AccountID string

type Account struct {
	ID            AccountID
	UserID        string
	AssetID       AssetID
	Balance       *big.Int
	LockedBalance *big.Int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserAccount struct {
	Account
	AssetSymbol string
}

func NewAccount(userID string, assetID AssetID) (*Account, error) {
	// TODO check assetID
	if userID == "" {
		return nil, ErrUserIDRequired
	}

	return &Account{
		UserID:        userID,
		AssetID:       assetID,
		Balance:       big.NewInt(0),
		LockedBalance: big.NewInt(0),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (a *Account) Deposit(amount *big.Int) error {
	if amount.Sign() <= 0 {
		return ErrInvalidAmount
	}
	a.Balance.Add(a.Balance, amount)
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Account) LockBalance(amount *big.Int) error {
	if amount.Sign() <= 0 {
		return ErrInvalidAmount
	}

	if a.Balance.Cmp(amount) < 0 {
		return ErrInSufficientBalance
	}

	a.Balance.Sub(a.Balance, amount)
	a.LockedBalance.Add(a.LockedBalance, amount)
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Account) SettleWithdrawal(amount *big.Int) error {
	if amount.Sign() <= 0 {
		return ErrInvalidAmount
	}
	if a.LockedBalance.Cmp(amount) < 0 {
		return ErrInsufficientLockedBalance
	}
	a.LockedBalance.Sub(a.LockedBalance, amount)
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Account) RevertWithdrawal(amount *big.Int) error {
	if amount.Sign() <= 0 {
		return ErrInvalidAmount
	}
	if a.LockedBalance.Cmp(amount) < 0 {
		return ErrInsufficientLockedBalance
	}

	a.LockedBalance.Sub(a.LockedBalance, amount)
	a.Balance.Add(a.Balance, amount)
	a.UpdatedAt = time.Now()
	return nil
}
