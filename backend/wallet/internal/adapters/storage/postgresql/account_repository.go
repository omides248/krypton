package postgresql

import (
	"context"
	"errors"
	"krypton/wallet/internal/domain"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type accountRepo struct {
	conn           *pgx.Conn
	defaultTimeout time.Duration
	logger         *zap.Logger
}

func NewAccountRepository(conn *pgx.Conn, timeout time.Duration, logger *zap.Logger) domain.AccountRepository {
	return &accountRepo{
		conn:           conn,
		defaultTimeout: timeout,
		logger:         logger.Named("postgres_account_repo"),
	}
}

func (r *accountRepo) Save(ctx context.Context, account *domain.Account) error {
	ctx, cancel := r.ctxWithTimeout(ctx)
	defer cancel()

	r.logger.Debug("saving a new account",
		zap.String("user_id", account.UserID),
		zap.String("asset_id", string(account.AssetID)),
	)

	const query = `
		INSERT INTO accounts (user_id, asset_id, balance, locked_balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.conn.QueryRow(ctx, query,
		account.UserID,
		account.AssetID,
		account.Balance.String(),
		account.LockedBalance.String(),
		account.CreatedAt,
		account.UpdatedAt,
	).Scan(&account.ID)

	if err != nil {
		r.logger.Error("failed to save account", zap.Error(err))
		return domain.ErrInternal
	}

	return nil
}

func (r *accountRepo) FindOrCreate(ctx context.Context, userID string, assetID domain.AssetID) (*domain.Account, error) {
	account, err := r.FindAccount(ctx, userID, assetID)
	if err == nil {
		return account, nil
	}

	// Internal err occur
	if !errors.Is(err, domain.ErrAccountNotFound) {
		return nil, err
	}

	r.logger.Debug("account not found, creating a new one",
		zap.String("user_id", userID),
		zap.String("asset_id", string(assetID)),
	)

	newAccount, err := domain.NewAccount(userID, assetID)
	if err != nil {
		return nil, err
	}

	if err := r.Save(ctx, newAccount); err != nil {
		return nil, err
	}

	return newAccount, nil
}

func (r *accountRepo) FindAccount(ctx context.Context, userID string, assetID domain.AssetID) (*domain.Account, error) {
	ctx, cancel := r.ctxWithTimeout(ctx)
	defer cancel()

	const query = `
			SELECT id, user_id, asset_id, balance, locked_balance, created_at, updated_at
			FROM accounts
			WHERE user_id = $1 AND asset_id = $2
    `

	var acc domain.Account
	var balanceStr, lockedBalanceStr string

	err := r.conn.QueryRow(ctx, query, userID, assetID).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.AssetID,
		&balanceStr,
		&lockedBalanceStr,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		r.logger.Error("failed to find account", zap.Error(err))
		return nil, domain.ErrInternal
	}

	// Convert string to decimal
	acc.Balance, _ = new(big.Int).SetString(balanceStr, 10)
	acc.LockedBalance, _ = new(big.Int).SetString(lockedBalanceStr, 10)

	return &acc, nil
}

func (r *accountRepo) FindUserAccounts(ctx context.Context, userID string) ([]domain.UserAccount, error) {
	ctx, cancel := r.ctxWithTimeout(ctx)
	defer cancel()

	const query = `
			SELECT a.id, a.user_id, a.asset_id, a.balance, a.locked_balance, a.created_at, a.updated_at, s.symbol
			FROM accounts a JOIN assets s ON a.asset_id = s.id
			WHERE a.user_id = $1
   `

	rows, err := r.conn.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error("failed to query user accounts", zap.Error(err))
		return nil, domain.ErrInternal
	}
	defer rows.Close()

	var userAccounts []domain.UserAccount
	for rows.Next() {
		var ua domain.UserAccount
		var balanceStr, lockedBalanceStr string
		if err := rows.Scan(&ua.ID, &ua.UserID, &ua.AssetID, &balanceStr, &lockedBalanceStr, &ua.CreatedAt, &ua.UpdatedAt, &ua.AssetSymbol); err != nil {
			r.logger.Error("failed to scan user account row", zap.Error(err))
			return nil, domain.ErrInternal
		}

		ua.Balance, _ = new(big.Int).SetString(balanceStr, 10)
		ua.LockedBalance, _ = new(big.Int).SetString(lockedBalanceStr, 10)
		userAccounts = append(userAccounts, ua)
	}

	return userAccounts, nil
}

func (r *accountRepo) ctxWithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, r.defaultTimeout)
}
