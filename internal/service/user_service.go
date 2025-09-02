package service

import (
	"context"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/MaksimPerv/Gofermart/internal/repository"
	"go.uber.org/zap"
)

type UserService interface {
	GetBalance(ctx context.Context, id int) (*entity.UserBalance, error)
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
