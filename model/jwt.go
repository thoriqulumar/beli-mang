package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Id       string `json:"id"`
	Username string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWTPayload struct {
	Id       string
	Username string
	Email    string
	Role     Role
}
