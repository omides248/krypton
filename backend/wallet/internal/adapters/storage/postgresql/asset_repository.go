package postgresql

import (
	"context"
	"krypton/wallet/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type assetRepo struct {
	conn           *pgx.Conn
	defaultTimeout time.Duration
	logger         *zap.Logger
}

func NewAssetRepository(conn *pgx.Conn, timeout time.Duration, logger *zap.Logger) (domain.AssetRepository, error) {
	return &assetRepo{
		conn:           conn,
		defaultTimeout: timeout,
		logger:         logger.Named("postgres_asset_repo"),
	}, nil
}

func (r *assetRepo) Save(ctx context.Context, asset *domain.Asset) error {
	ctx, cancel := r.ctxWithTimeout(ctx)
	defer cancel()

	r.logger.Debug("saving a new asset", zap.String("symbol", asset.Symbol))

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return domain.ErrInternal
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	const assetQuery = `
					INSERT INTO assets (symbol, name, asset_type, precision, is_active, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
					RETURNING id
	`
	err = tx.QueryRow(ctx, assetQuery,
		asset.Symbol, asset.Name, asset.Type, asset.Precision, asset.IsActive, asset.CreatedAt, asset.UpdatedAt,
	).Scan(&asset.ID)
	if err != nil {
		r.logger.Error("failed to save asset", zap.Error(err))
		return domain.ErrInternal
	}

	for _, an := range asset.SupportedNetworks {
		const networkQuery = `
			INSERT INTO asset_networks (asset_id, network_id,  contract_address, min_withdrawal, withdrawal_fee, is_active)
			VALUES ($1, $2, $3, $4, $5, $6)
       `
		_, err = tx.Exec(ctx, networkQuery,
			asset.ID, an.NetworkID, an.ContractAddress, an.MinWithdrawal, an.WithdrawalFee, an.IsActive,
		)
		if err != nil {
			r.logger.Error("failed to save asset_network link", zap.Error(err))
			return domain.ErrInternal
		}
	}

	return tx.Commit(ctx)
}

func (r *assetRepo) FindAssetBySymbol(ctx context.Context, symbol string) (*domain.Asset, error) {
	//TODO implement me
	panic("implement me")
}

func (r *assetRepo) FindUserAccounts(ctx context.Context, userID string) ([]domain.UserAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (r *assetRepo) ctxWithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, r.defaultTimeout)
}
