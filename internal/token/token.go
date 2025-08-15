package token

import (
	"errors"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set in environment")
	}
	return []byte(secret), nil
}

func GenerateToken(user *entity.User) (string, error) {
	jwtSecret, err := getJWTSecret()
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := Claims{
		Username: user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
