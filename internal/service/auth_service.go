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
	Register(ctx context.Context, user *entity.User) (int, error)
	GenerateToken(user *entity.User, id *int) (string, error)
	Login(ctx context.Context, user *entity.User) (int, error)
	ValidateToken(tokenString string) (*token.Claims, error)
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

func (a *authService) Register(ctx context.Context, user *entity.User) (int, error) {
	if _, err := a.repo.GetByUser(ctx, user.Login); err == nil {
		return 0, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(hashedPassword)
	return a.repo.CreateUser(ctx, user)
}

func (a *authService) GenerateToken(user *entity.User, id *int) (string, error) {
	return token.GenerateToken(user, id)
}
func (a *authService) Login(ctx context.Context, user *entity.User) (int, error) {
	consumer, err := a.repo.GetByUser(ctx, user.Login)
	if err != nil {
		a.logger.Error("Failed to get user from database", zap.String("username", user.Login), zap.Error(err))
		return 0, err
	}
	if consumer == nil {
		return 0, ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(consumer.Password), []byte(user.Password)); err != nil {
		return 0, ErrInvalidCredentials
	}
	return consumer.Id, nil
}

func (a *authService) ValidateToken(tokenString string) (*token.Claims, error) {
	return token.ValidateToken(tokenString)
}
