package domain

import (
	"math/big"
	"time"
)

type AssetID string

type AssetType string

const (
	AssetTypeCrypto AssetType = "CRYPTO"
	AssetTypeFiat   AssetType = "FIAT"
)

type Asset struct {
	ID                AssetID
	Symbol            string
	Name              string
	Type              AssetType
	Precision         int
	MinWithdrawal     *big.Int
	WithdrawalFee     *big.Int
	IsActive          bool
	SupportedNetworks []AssetNetwork
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type AssetNetwork struct {
	AssetID         AssetID
	NetworkID       NetworkID
	ContractAddress *string // Exp: ERC-20 token
	MinWithdrawal   *big.Int
	WithdrawalFee   *big.Int
	IsActive        bool
}

func NewAsset(symbol, name, network string, assetType AssetType, precision int) (*Asset, error) {
	if symbol == "" || name == "" || network == "" {
		return nil, ErrAssetFieldsRequired
	}

	return &Asset{
		Symbol:    symbol,
		Name:      name,
		Type:      assetType,
		Precision: precision,
		IsActive:  true,
		CreatedAt: time.Now(),
	}, nil
}
