package postgresql

import (
	"context"
	"krypton/wallet/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type networkRepo struct {
	conn           *pgx.Conn
	defaultTimeout time.Duration
	logger         *zap.Logger
}

func NewNetworkRepository(conn *pgx.Conn, defaultTimeout time.Duration, logger *zap.Logger) domain.NetworkRepository {
	return &networkRepo{
		conn:           conn,
		defaultTimeout: defaultTimeout,
		logger:         logger.Named("postgres_network_repo"),
	}
}

func (r *networkRepo) Save(ctx context.Context, network *domain.Network) error {
	ctx, cancel := r.ctxWithTimeout(ctx)
	defer cancel()

	r.logger.Debug("saving a new network", zap.String("name", network.Name))

	const query = `
				INSERT INTO networks (name, chain_id, rpc_url, is_active, created_at)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id;
	`

	err := r.conn.QueryRow(ctx, query,
		network.Name, network.ChainID, network.RpcURL, network.IsActive, network.CreatedAt,
	).Scan(&network.ID)

	if err != nil {
		r.logger.Error("failed to save network", zap.Error(err))
		return domain.ErrInternal
	}

	return nil
}

func (r *networkRepo) FindOrCreate(ctx context.Context, userID string, assetID domain.AssetID) (*domain.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (r *networkRepo) ctxWithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, r.defaultTimeout)
}
