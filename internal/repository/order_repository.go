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
	GetOrders(context.Context, int) ([]entity.DBUser, error)
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
	rows, err := p.db.Query(ctx, "SELECT od.number,os.name,od.accrual,od.created_at FROM orders od JOIN order_statuses os on os.id= od.status_id WHERE od.user_id=$1 ORDER BY od.created_at", userId)
	if err != nil {
		p.logger.Error("Request execution error", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	var users []entity.DBUser

	for rows.Next() {
		var user entity.DBUser
		if err = rows.Scan(&user.Number, &user.Status, &user.Accrual, &user.UploadedAt); err != nil {
			p.logger.Error("Line scan error", zap.Error(err))
			return nil, err
		}
		if user.Accrual.IsZero() {
			user.Accrual = nil
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		p.logger.Error("Error while iterating over rows", zap.Error(err))
		return nil, err
	}
	return users, nil
}
