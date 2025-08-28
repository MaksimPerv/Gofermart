package service

import (
	"github.com/MaksimPerv/Gofermart/internal/repository"
	"go.uber.org/zap"
)

type OrderService interface {
	CreateOrder(string) error
}

type orderService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewOrderService(repo repository.UserRepository, logger *zap.Logger) OrderService {
	return &orderService{
		repo:   repo,
		logger: logger,
	}
}

func (o *orderService) CreateOrder(number string) error {
	return nil
}
