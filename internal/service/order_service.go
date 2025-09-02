package service

import (
	"context"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/MaksimPerv/Gofermart/internal/repository"
	"go.uber.org/zap"
)

type OrderService interface {
	CreateOrder(context.Context, string, int) error
	GetOrders(context.Context, int) ([]entity.DBUser, error)
}

type orderService struct {
	repo   repository.OrderRepository
	logger *zap.Logger
}

func NewOrderService(repo repository.OrderRepository, logger *zap.Logger) OrderService {
	return &orderService{
		repo:   repo,
		logger: logger,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, number string, userID int) error {
	o.logger.Debug("Creating order",
		zap.String("order_number", number),
		zap.Int("user_id", userID))

	ok, err := o.repo.CreateOrder(ctx, number, userID)
	if ok {
		o.logger.Info("Order created successfully",
			zap.String("order_number", number),
			zap.Int("user_id", userID))
		return nil
	}

	if err == nil {
		o.logger.Debug("Order already exists,checking owner",
			zap.String("order_number", number))

		id, check := o.repo.GetID(ctx, number)
		if check != nil {
			o.logger.Error("Failed to get order owner",
				zap.String("order_number", number),
				zap.Error(check))
			return check
		}
		if id == userID {
			o.logger.Info("Order already uploaded by same user",
				zap.String("order_number", number),
				zap.Int("user_id", userID))
			return ErrOrderAlreadyUploaded
		} else {
			o.logger.Info("Order belongs to another user",
				zap.String("order_number", number),
				zap.Int("requested_by", userID),
				zap.Int("owned_by", id))
			return ErrOrderBelongsToAnotherUser
		}
	}
	o.logger.Error("Failed to create order",
		zap.String("order_number", number),
		zap.Int("user_id", userID),
		zap.Error(err))
	return err

}

func (o *orderService) GetOrders(ctx context.Context, userId int) ([]entity.DBUser, error) {
	return o.repo.GetOrders(ctx, userId)
}
