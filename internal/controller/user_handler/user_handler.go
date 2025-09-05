package user_handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/MaksimPerv/Gofermart/internal/service"
	"github.com/MaksimPerv/Gofermart/internal/validation"
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

func (u *UserHandler) WithdrawalRequest(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("CreateOrder request started",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path))

	if r.Header.Get("Content-Type") != "application/json" {
		u.logger.Warn("Invalid Content-Type",
			zap.String("content_type", r.Header.Get("Content-Type")))
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		u.logger.Error("User not authenticated in context")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()
	var req entity.WithdrawRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		u.logger.Error("Failed to marshal users to JSON",
			zap.Int("user_id", userId),
			zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !validation.ValidateLuhn(req.Order) {
		u.logger.Warn("Invalid order number format",
			zap.String("order_number", req.Order))
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}
	err = u.userService.WithdrawRequest(context.Background(), req.Sum, userId, req.Order)

	if errors.Is(err, service.ErrInsufficientPoints) {
		http.Error(w, "there are insufficient funds in the account", http.StatusPaymentRequired)
		return
	}
	if err != nil {
		u.logger.Error("Failed to process withdrawal request",
			zap.Int("user_id", userId),
			zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successful request processing;"))
}

func (u *UserHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("GetWithdrawals request started",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path))

	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		u.logger.Error("User not authenticated in context")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	response, err := u.userService.GetWithdrawals(context.Background(), userId)
	if err != nil {
		u.logger.Error("Error get withdrawals", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if response == nil {
		u.logger.Info("No information to answer",
			zap.Int("user", userId))
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("No information to answer"))
		return
	}

	withdrawalJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		u.logger.Error("Failed to marshal users to JSON",
			zap.Int("user_id", userId),
			zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(withdrawalJSON)
	if err != nil {
		u.logger.Error("Failed to write JSON response",
			zap.Int("user_id", userId),
			zap.Error(err),
		)
	}
}
