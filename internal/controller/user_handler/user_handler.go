package user_handler

import (
	"github.com/MaksimPerv/Gofermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	logger      *zap.Logger
	userService service.UserService
}

func NewUserHandler(logger *zap.Logger, userService service.UserService) *UserHandler {
	return &UserHandler{logger: logger, userService: userService}
}

func (u *UserHandler) Get(w http.ResponseWriter, r *http.Request) {

}
