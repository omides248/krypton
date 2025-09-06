package adapters

import (
	"fmt"
	"krypton/identity/internal/adapters/storage/postgresql"
	"krypton/identity/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Adapter struct {
	UserRepo domain.UserRepository
}

func NewAdapter(db *pgx.Conn, logger *zap.Logger) (*Adapter, error) {
	dbTimeout := 10 * time.Second
	userRepo, err := postgresql.NewUserRepository(db, logger, dbTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create product repository: %w", err)
	}

	return &Adapter{
		UserRepo: userRepo,
	}, nil

}
