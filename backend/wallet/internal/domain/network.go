package domain

import "time"

type NetworkID string

type Network struct {
	ID        NetworkID
	Name      string
	ChainID   *int64 // Empty not network is not EVM
	RpcURL    string
	IsActive  bool
	CreatedAt time.Time
}

func NewNetwork(name, rpcURL string, chainID *int64) (*Network, error) {
	if name == "" || rpcURL == "" {
		return nil, ErrNetworkFieldsRequired
	}

	return &Network{
		Name:      name,
		ChainID:   chainID,
		RpcURL:    rpcURL,
		IsActive:  true,
		CreatedAt: time.Now(),
	}, nil
}
