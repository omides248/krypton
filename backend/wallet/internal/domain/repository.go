package domain

import "context"

type AssetRepository interface {
	Save(ctx context.Context, asset *Asset) error
	FindAssetBySymbol(ctx context.Context, symbol string) (*Asset, error)
}

type NetworkRepository interface {
	Save(ctx context.Context, network *Network) error
}

type AccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindAccount(ctx context.Context, userID string, assetID AssetID) (*Account, error)
}
