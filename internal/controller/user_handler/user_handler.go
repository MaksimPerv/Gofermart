package user_handler

import (
	"context"
	"encoding/json"
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

func (u *UserHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("GetBalance request started",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path))

	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		u.logger.Error("User not authenticated in context")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	u.logger.Info("Get balance in service",
		zap.Int("user_id", userId))

	balance, err := u.userService.GetBalance(context.Background(), userId)
	if err != nil {
		u.logger.Error("Error get balance",
			zap.Int("user_id", userId),
			zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		u.logger.Error("Failed to marshal users to JSON",
			zap.Int("user_id", userId),
			zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(balanceJSON)
	if err != nil {
		u.logger.Error("Failed to write JSON response",
			zap.Int("user_id", userId),
			zap.Error(err),
		)
	}
}
