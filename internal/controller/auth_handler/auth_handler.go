package auth_handler

import (
	"encoding/json"
	"errors"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/MaksimPerv/Gofermart/internal/service"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AuthHandler struct {
	logger      *zap.Logger
	authService service.AuthService
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
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("Invalid JSON", zap.Error(err))
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}
	userId, err := a.authService.Register(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			http.Error(w, "Login already taken", http.StatusConflict)
		default:
			a.logger.Error("Registration failed", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	token, err := a.authService.GenerateToken(&req, &userId)
	if err != nil {
		a.logger.Error("Token generation failed", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered and authenticated"))
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var user entity.User
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		a.logger.Error("Invalid JSON", zap.Error(err))
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if user.Login == "" || user.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}
	userId, err := a.authService.Login(r.Context(), &user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		default:
			a.logger.Error("Service error in login", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	token, err := a.authService.GenerateToken(&user, &userId)
	if err != nil {
		a.logger.Error("Token generation failed", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	a.logger.Info("Login successful", zap.String("username", user.Login))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully authenticated"))
}
