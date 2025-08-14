package repository

import (
	"context"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
}

type postgresUserRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresUserRepository(db *pgxpool.Pool, logger *zap.Logger) UserRepository {
	return &postgresUserRepository{
		db:     db,
		logger: logger,
	}
}

func (postgre *postgresUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return nil, nil
}
