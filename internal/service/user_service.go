package service

import (
	"context"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/MaksimPerv/Gofermart/internal/repository"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type UserService interface {
	GetBalance(ctx context.Context, id int) (*entity.UserBalance, error)
	WithdrawRequest(ctx context.Context, sum decimal.Decimal, userId int) error
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (u *userService) GetBalance(ctx context.Context, id int) (*entity.UserBalance, error) {
	return u.repo.GetBalance(ctx, id)
}

func (u *userService) WithdrawRequest(ctx context.Context, sum decimal.Decimal, userId int) error {
	return u.repo.ProcessWithdrawal(ctx, sum, userId)
}
