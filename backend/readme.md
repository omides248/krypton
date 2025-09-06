identity-service

(POST) /v1/auth/register
(POST) /v1/auth/login
(GET) /v1/users/me

--------------------------------------
wallet-service

Accounts
(GET) /v1/accounts
(GET) /v1/accounts/{asset_symbol}

Deposit Addresses
(POST) /v1/addresses
(GET) /v1/addresses/{asset_symbol}

Withdrawals
(POST) /v1/withdrawals
(GET) /v1/withdrawals

Transactions
(GET) /v1/transactions
--------------------------------------
blockchain-scanner-service