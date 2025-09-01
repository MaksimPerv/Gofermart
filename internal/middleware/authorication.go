package middleware

import (
	"context"
	"github.com/MaksimPerv/Gofermart/internal/service"
	"net/http"
	"time"
)

func AuthMiddleware(a service.AuthService) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			token := cookie.Value
			claims, err := a.ValidateToken(token)
			if err != nil {
				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Value:    "",
					Path:     "/",
					Expires:  time.Now().Add(-time.Hour),
					HttpOnly: true,
					Secure:   false,
					SameSite: http.SameSiteLaxMode,
				})
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
