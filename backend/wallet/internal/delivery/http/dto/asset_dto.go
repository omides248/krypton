package dto

import (
	"krypton/wallet/internal/domain"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateAssetRequest struct {
	Name      string           `json:"name"`
	Symbol    string           `json:"symbol"`
	Type      domain.AssetType `json:"type"`
	Precision int              `json:"precision"`
}

func (r CreateAssetRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Name, ozzo.Required, ozzo.Length(2, 50)),
		ozzo.Field(&r.Symbol, ozzo.Required, ozzo.Length(2, 10)),
		ozzo.Field(&r.Type, ozzo.Required, ozzo.In(domain.AssetTypeCrypto, domain.AssetTypeFiat)),
		ozzo.Field(&r.Precision, ozzo.Required, ozzo.Min(0), ozzo.Max(18)),
	)
}

type AssetResponse struct {
	ID        string           `json:"id"`
	Symbol    string           `json:"symbol"`
	Name      string           `json:"name"`
	Type      domain.AssetType `json:"type"`
	Precision int              `json:"precision"`
	IsActive  bool             `json:"is_active"`
}

func ToAssetResponse(asset *domain.Asset) AssetResponse {
	return AssetResponse{
		ID:        string(asset.ID),
		Symbol:    asset.Symbol,
		Name:      asset.Name,
		Type:      asset.Type,
		Precision: asset.Precision,
		IsActive:  asset.IsActive,
	}
}

func ToAssetsResponse(assets []*domain.Asset) []AssetResponse {
	response := make([]AssetResponse, len(assets))
	for i, asset := range assets {
		response[i] = ToAssetResponse(asset)
	}
	return response
}

type CreateNetworkRequest struct {
	Name    string `json:"name"`
	ChainID *int64 `json:"chain_id"`
	RpcURL  string `json:"rpc_url"`
}

func (r CreateNetworkRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Name, ozzo.Required, ozzo.Length(2, 50)),
		ozzo.Field(&r.RpcURL, ozzo.Required, ozzo.Length(5, 255)),
	)
}

type NetworkResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ChainID  *int64 `json:"chain_id,omitempty"`
	RpcURL   string `json:"rpc_url"`
	IsActive bool   `json:"is_active"`
}

func ToNetworkResponse(network *domain.Network) NetworkResponse {
	return NetworkResponse{
		ID:       string(network.ID),
		Name:     network.Name,
		ChainID:  network.ChainID,
		RpcURL:   network.RpcURL,
		IsActive: network.IsActive,
	}
}

func ToNetworksResponse(networks []*domain.Network) []NetworkResponse {
	response := make([]NetworkResponse, len(networks))
	for i, network := range networks {
		response[i] = ToNetworkResponse(network)
	}
	return response
}

type AddAssetToNetworkRequest struct {
	NetworkName     string `json:"network_name"`
	ContractAddress string `json:"contract_address"` // optional
	MinWithdrawal   string `json:"min_withdrawal"`
	WithdrawalFee   string `json:"withdrawal_fee"`
}

func (r AddAssetToNetworkRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.NetworkName, ozzo.Required),
		ozzo.Field(&r.MinWithdrawal, ozzo.Required),
		ozzo.Field(&r.WithdrawalFee, ozzo.Required),
	)
}
