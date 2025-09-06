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

func NewAssetRepository(conn *pgx.Conn, logger *zap.Logger, timeout time.Duration) (domain.AssetRepository, error) {
	return &assetRepo{
		conn:           conn,
		defaultTimeout: timeout,
		logger:         logger.Named("postgres_asset_repo"),
	}, nil
}

func (a assetRepo) Save(ctx context.Context, asset *domain.Asset) error {
	//TODO implement me
	panic("implement me")
}

func (a assetRepo) FindAssetBySymbol(ctx context.Context, symbol string) (*domain.Asset, error) {
	//TODO implement me
	panic("implement me")
}
