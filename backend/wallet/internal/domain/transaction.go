package domain

import (
	"math/big"
	"time"
)

// --- Deposit ---

type DepositID string
type DepositStatus int

const (
	DepositStatusUnconfirmed DepositStatus = iota
	DepositStatusConfirmed
)

type Deposit struct {
	ID            DepositID
	AccountID     AccountID
	AssetID       AssetID
	Status        DepositStatus
	Amount        *big.Int
	OnchainTxID   string
	FromAddress   string
	Confirmations int
	CreatedAt     time.Time
	ConfirmedAt   *time.Time
}

// --- Withdrawal ---

type WithdrawalID string
type WithdrawalStatus int

const (
	WithdrawalStatusRequested WithdrawalStatus = iota
	WithdrawalStatusApproved
	WithdrawalStatusProcessing
	WithdrawalStatusBroadcasted
	WithdrawalStatusConfirmed
	WithdrawalStatusFailed
)

type Withdrawal struct {
	ID          WithdrawalID
	AccountID   AccountID
	AssetID     AssetID
	Status      WithdrawalStatus
	Amount      *big.Int
	ToAddress   string
	OnchainTxID *string
	NetworkFee  *big.Int
	ApprovedBy  *string
	Priority    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InternalTransferID string

type InternalTransfer struct {
	ID            InternalTransferID
	FromAccountID AccountID
	ToAccountID   AccountID
	AssetID       AssetID
	Amount        *big.Int
	ReferenceID   *string
	CreatedAt     time.Time
}

func NewInternalTransfer(from, to AccountID, assetID AssetID, amount *big.Int, referenceID *string) (*InternalTransfer, error) {
	if from == to {
		return nil, ErrCannotTransferToSameAccount
	}
	if amount.Sign() <= 0 {
		return nil, ErrInvalidAmount
	}

	return &InternalTransfer{
		FromAccountID: from,
		ToAccountID:   to,
		AssetID:       assetID,
		Amount:        amount,
		ReferenceID:   referenceID,
		CreatedAt:     time.Now(),
	}, nil
}
