package dto

import "krypton/wallet/internal/domain"

type AccountResponse struct {
	AssetSymbol   string `json:"asset_symbol"`
	Balance       string `json:"balance"`
	LockedBalance string `json:"locked_balance"`
}

func ToAccountResponse(userAccount domain.UserAccount) AccountResponse {
	return AccountResponse{
		AssetSymbol:   userAccount.AssetSymbol,
		Balance:       userAccount.Balance.String(),
		LockedBalance: userAccount.LockedBalance.String(),
	}
}

func ToAccountsResponse(accounts []domain.UserAccount) []AccountResponse {
	responses := make([]AccountResponse, len(accounts))
	for i, acc := range accounts {
		responses[i] = ToAccountResponse(acc)
	}

	return responses
}
