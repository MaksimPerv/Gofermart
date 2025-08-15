package token

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Username             string `json:"username"`
	jwt.RegisteredClaims        // стандартные поля (exp, iss и др.)
}
