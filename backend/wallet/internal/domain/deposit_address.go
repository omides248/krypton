package domain

import (
	"time"
)

type DepositAddressID string

type DepositAddress struct {
	ID             DepositAddressID
	AccountID      AccountID
	Address        string
	DerivationPath string
	MemoTag        *string
	CreatedAt      time.Time
}

func NewDepositAddress(accountID AccountID, address, derivationPath string, memoTag *string) (*DepositAddress, error) {
	if address == "" {
		return nil, ErrAddressRequired
	}

	return &DepositAddress{
		AccountID:      accountID,
		Address:        address,
		DerivationPath: derivationPath,
		MemoTag:        memoTag,
		CreatedAt:      time.Now(),
	}, nil
}
