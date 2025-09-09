package domain

import "context"

type AssetRepository interface {
	Save(ctx context.Context, asset *Asset) error
	AddAssetToNetwork(ctx context.Context, assetNetwork *AssetNetwork) error
	FindAssetBySymbol(ctx context.Context, symbol string) (*Asset, error)
	FindAllAssets(ctx context.Context) ([]*Asset, error)
}

type NetworkRepository interface {
	Save(ctx context.Context, network *Network) error
	FindOrCreate(ctx context.Context, userID string, assetID AssetID) (*Account, error)
	FindAllNetworks(ctx context.Context) ([]*Network, error)
}

type AccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindOrCreate(ctx context.Context, userID string, assetID AssetID) (*Account, error)
	FindAccount(ctx context.Context, userID string, assetID AssetID) (*Account, error)
	FindUserAccounts(ctx context.Context, userID string) ([]UserAccount, error)
	//Update(ctx context.Context, account *Account) error
}

type DepositAddressRepository interface {
	Save(ctx context.Context, addr *DepositAddress) error
	FindByAccountID(ctx context.Context, accountID AccountID) (*DepositAddress, error)
}
