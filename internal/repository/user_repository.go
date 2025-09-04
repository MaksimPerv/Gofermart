package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetByUser(ctx context.Context, login string) (*entity.UserWithId, error)
	CreateUser(ctx context.Context, user *entity.User) (int, error)
	GetBalance(context.Context, int) (*entity.UserBalance, error)
	ProcessWithdrawal(ctx context.Context, amount decimal.Decimal, userID int) error
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
func (p *postgresUserRepository) ProcessWithdrawal(ctx context.Context, amount decimal.Decimal, userID int) error {
	tx, err := p.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		p.logger.Error("Error start transaction",
			zap.Error(err))
		return err
	}
	var balance decimal.Decimal

	err = tx.QueryRow(ctx, "SELECT current  from balance WHERE user_id=$1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		p.logger.Error("Error from select balance",
			zap.Int("user_id", userID),
			zap.Error(err))
		return err
	}
	if !balance.GreaterThanOrEqual(amount) {
		return errors.New("insufficient points")
	}

	_, err = tx.Exec(ctx, "UPDATE balance SET current=current-$1,withdrawn=withdrawn+$1 WHERE user_id=$2", amount, userID)
	if err != nil {
		p.logger.Error("Error update transaction",
			zap.Int("user_id", userID),
			zap.Error(err))
		return err
	}
	return tx.Commit(ctx)
}
