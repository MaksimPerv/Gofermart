package service

import (
	"context"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/MaksimPerv/Gofermart/internal/repository"
	"github.com/MaksimPerv/Gofermart/internal/token"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, user *entity.User) error
	GenerateToken(user *entity.User) (string, error)
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

func (a *authService) Register(ctx context.Context, user *entity.User) error {
	if _, err := a.repo.GetByUser(ctx, user.Login); err == nil {
		return ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return a.repo.CreateUser(ctx, user)
}

func (a *authService) GenerateToken(user *entity.User) (string, error) {
	return token.GenerateToken(user)
}
