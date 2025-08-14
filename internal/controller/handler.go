package controller

import (
	"github.com/MaksimPerv/Gofermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	logger      *zap.Logger
	userService service.UserService
}

type AuthHandler struct {
	logger      *zap.Logger
	AuthService service.AuthService
}

func NewUserHandler(logger *zap.Logger, userService service.UserService) *UserHandler {
	return &UserHandler{logger: logger, userService: userService}
}

func NewAuthHandler(logger *zap.Logger, AuthService service.AuthService) *AuthHandler {
	return &AuthHandler{logger: logger, AuthService: AuthService}
}

func (u *UserHandler) Get(w http.ResponseWriter, r *http.Request) {

}
func (u *AuthHandler) Get(w http.ResponseWriter, r *http.Request) {

}
