package postgresql

import (
	"context"
	"errors"
	"krypton/identity/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type userRepo struct {
	conn           *pgx.Conn
	defaultTimeout time.Duration
	logger         *zap.Logger
}

func NewUserRepository(conn *pgx.Conn, logger *zap.Logger, timeout time.Duration) (domain.UserRepository, error) {
	return &userRepo{
		conn:           conn,
		defaultTimeout: timeout,
		logger:         logger.Named("postgres_user_repo"),
	}, nil
}

func (r *userRepo) Save(ctx context.Context, user *domain.User) error {
	ctx, cancel := r.ctxWithTimeout(ctx)
	defer cancel()

	if !user.Status.IsValid() {
		r.logger.Warn("attempt to save user with invalid status",
			zap.String("email", user.Email),
			zap.Int("status", int(user.Status)),
		)
		return domain.ErrInvalidUserStatus
	}

	r.logger.Debug("saving a new user", zap.String("email", user.Email))

	const query = `
			INSERT INTO users (email, password_hash, status, kyc_level, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
	`

	err := r.conn.QueryRow(ctx, query,
		user.Email, user.PasswordHash, user.Status, user.KYCLevel, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			r.logger.Error("no rows on insert returning", zap.Error(err))
			return domain.ErrInternal
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			r.logger.Info("email already exists", zap.String("email", user.Email))
			return domain.ErrEmailAlreadyExists
		default:
			r.logger.Error("failed to save user", zap.String("email", user.Email), zap.Error(err))
			return domain.ErrInternal
		}
	}

	return nil
}

func (r *userRepo) find(ctx context.Context, query string, args ...interface{}) (*domain.User, error) {
	ctx, cancel := r.ctxWithTimeout(ctx)
	defer cancel()

	var user domain.User

	err := r.conn.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Status,
		&user.KYCLevel,
		&user.TwoFactorSecret,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("failed to query user", zap.Any("args", args), zap.Error(err))
		return nil, domain.ErrInternal
	}

	if !user.Status.IsValid() {
		r.logger.Error("invalid user status from DB",
			zap.String("user_id", string(user.ID)),
			zap.Int("status", int(user.Status)),
		)
		return nil, domain.ErrInvalidUserStatus
	}

	return &user, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.logger.Debug("finding user by email", zap.String("email", email))

	const query = `
			SELECT id, email, password_hash, status, kyc_level, two_factor_secret, created_at, updated_at 
			FROM users 
			WHERE email = $1
	`

	user, err := r.find(ctx, query, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		r.logger.Error("failed to find user by email", zap.String("email", email), zap.Error(err))
	}

	return user, err
}

func (r *userRepo) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	r.logger.Debug("finding user by id", zap.String("user_id", string(id)))

	const query = `
			SELECT id, email, password_hash, status,kyc_level, two_factor_secret, created_at, updated_at 
			FROM users 
			WHERE id = $1
	`

	user, err := r.find(ctx, query, id)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		r.logger.Error("failed to find user by id", zap.String("user_id", string(id)), zap.Error(err))
	}

	return user, nil
}

func (r *userRepo) ctxWithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, r.defaultTimeout)
}
