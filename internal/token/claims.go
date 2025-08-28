package token

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
