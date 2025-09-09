package domain

import "context"

type AssetRepository interface {
	Save(ctx context.Context, asset *Asset) error
	FindAssetBySymbol(ctx context.Context, symbol string) (*Asset, error)
	//SaveAssetNetwork(ctx context.Context, assetNetwork *AssetNetwork) error
}

type NetworkRepository interface {
	Save(ctx context.Context, network *Network) error
	FindOrCreate(ctx context.Context, userID string, assetID AssetID) (*Account, error)
	//FindByName(ctx context.Context, name string) (*Network, error)
}

type AccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindOrCreate(ctx context.Context, userID string, assetID AssetID) (*Account, error)
	FindAccount(ctx context.Context, userID string, assetID AssetID) (*Account, error)
	FindUserAccounts(ctx context.Context, userID string) ([]UserAccount, error)
	//UpdateBalance(ctx context.Context, account *Account) error
	//FindByID(ctx context.Context, id AccountID) (*Account, error)
}
