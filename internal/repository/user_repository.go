package repository

import (
	"context"
	"database/sql"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetByUser(ctx context.Context, login string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
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

func (postgre *postgresUserRepository) GetByUser(ctx context.Context, login string) (*entity.User, error) {
	var user entity.User
	err := postgre.db.QueryRow(ctx, "SELECT login,password FROM users WHERE login=$1;", login).Scan(&user.Login, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (postgre *postgresUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	_, err := postgre.db.Exec(ctx, "INSERT INTO users (login,password) VALUES ($1,$2);", user.Login, user.Password)
	return err
}
