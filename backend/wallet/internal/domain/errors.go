package domain

import "errors"

var (
	ErrAssetNotFound              = errors.New("asset not found")
	ErrAssetNotSupportedOnNetwork = errors.New("asset is not supported on the specified network")
	ErrAssetFieldsRequired        = errors.New("asset symbol and name are required")

	ErrNetworkNotFound       = errors.New("network not found")
	ErrNetworkFieldsRequired = errors.New("network name and rpc_url are required")

	ErrAccountNotFound = errors.New("account not found")
	ErrUserIDRequired  = errors.New("user id cannot be empty")

	ErrInvalidAmount             = errors.New("amount must be positive")
	ErrInsufficientLockedBalance = errors.New("insufficient locked balance")
	ErrInSufficientBalance       = errors.New("insufficient balance")

	ErrDepositNotFound = errors.New("deposit not found")

	ErrWithdrawalNotFound        = errors.New("withdrawal not found")
	ErrWithdrawalAmountTooLow    = errors.New("withdrawal amount is below the minimum limit")
	ErrWithdrawalAddressRequired = errors.New("withdrawal address cannot be empty")

	ErrDeposit = errors.New("cannot transfer to the same account")

	ErrAddressRequired = errors.New("address cannot be empty")

	ErrCannotTransferToSameAccount = errors.New("cannot transfer to the same account")

	ErrInternal = errors.New("internal server error")
)
