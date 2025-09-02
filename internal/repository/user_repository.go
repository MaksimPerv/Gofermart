package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetByUser(ctx context.Context, login string) (*entity.UserWithId, error)
	CreateUser(ctx context.Context, user *entity.User) (int, error)
	GetBalance(context.Context, int) (*entity.UserBalance, error)
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

func (postgre *postgresUserRepository) GetByUser(ctx context.Context, login string) (*entity.UserWithId, error) {
	var user entity.UserWithId
	err := postgre.db.QueryRow(ctx, "SELECT id,login,password FROM users WHERE login=$1;", login).Scan(&user.Id, &user.Login, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (postgre *postgresUserRepository) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	var userId int
	err := postgre.db.QueryRow(ctx, "INSERT INTO users (login,password) VALUES ($1,$2) RETURNING id;", user.Login, user.Password).Scan(&user)
	return userId, err
}

func (p *postgresUserRepository) GetBalance(ctx context.Context, userId int) (*entity.UserBalance, error) {
	var userBalance entity.UserBalance
	err := p.db.QueryRow(ctx, "SELECT current,withdrawn FROM balance WHERE user_id=$1", userId).Scan(&userBalance.Current, &userBalance.Withdrawn)
	if err != nil {
		return nil, err
	}
	return &userBalance, nil
}
