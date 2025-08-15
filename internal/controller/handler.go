package controller

import (
	"encoding/json"
	"errors"
	"github.com/MaksimPerv/Gofermart/internal/entity"
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
	authService service.AuthService
}

func NewUserHandler(logger *zap.Logger, userService service.UserService) *UserHandler {
	return &UserHandler{logger: logger, userService: userService}
}

func NewAuthHandler(logger *zap.Logger, AuthService service.AuthService) *AuthHandler {
	return &AuthHandler{logger: logger, authService: AuthService}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var req entity.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("Invalid JSON", zap.Error(err))
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	if err := a.authService.Register(r.Context(), &req); err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			http.Error(w, "Login already taken", http.StatusConflict)
		default:
			a.logger.Error("Registration failed", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	token, err := a.authService.GenerateToken(&req)
	if err != nil {
		a.logger.Error("Token generation failed", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:        "token",
		Value:       token,
		Path:        "/",
		HttpOnly:    true,
		Partitioned: false,
		Raw:         "",
		Unparsed:    nil,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered and authenticated"))
}

func (u *UserHandler) Get(w http.ResponseWriter, r *http.Request) {

}
func (a *AuthHandler) Get(w http.ResponseWriter, r *http.Request) {

}
