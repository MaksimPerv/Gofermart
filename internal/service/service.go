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

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

type authService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewAuthService(repo repository.UserRepository, logger *zap.Logger) AuthService {
	return &authService{
		repo:   repo,
		logger: logger,
	}
}

func (a *authService) Login(ctx context.Context, email, password string) (string, error) {
	return "nil", nil
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
