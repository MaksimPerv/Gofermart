package repository

import (
	"context"
	"errors"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type OrderRepository interface {
	CreateOrder(context.Context, string, int) (bool, error)
	GetID(ctx context.Context, number string) (int, error)
}
type postgresOrderRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresOrderRepository(db *pgxpool.Pool, logger *zap.Logger) OrderRepository {
	return &postgresOrderRepository{
		db:     db,
		logger: logger,
	}
}

func (o *postgresOrderRepository) CreateOrder(ctx context.Context, number string, userID int) (bool, error) {
	_, err := o.db.Exec(ctx, "INSERT INTO orders (number,user_id) VALUES($1,$2)", number, userID)

	if err == nil {
		return true, nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return false, nil
	}

	return false, err
}

func (o *postgresOrderRepository) GetID(ctx context.Context, number string) (int, error) {
	var id int
	err := o.db.QueryRow(ctx, "SELECT id FROM orders WHERE number=$1", number).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *postgresOrderRepository) GetOrders(ctx context.Context, userId int) ([]entity.DBUser, error) {
	return nil, nil
}
