package service

import (
	"context"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/MaksimPerv/Gofermart/internal/repository"
	"go.uber.org/zap"
)

type UserService interface {
	Get(ctx context.Context, id string) (*entity.User, error)
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

func (u *userService) Get(ctx context.Context, id string) (*entity.User, error) {
	return nil, nil
}
