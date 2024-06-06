package crypto

import (
	"beli-mang/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(staff model.Staff, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		Id:       staff.ID.String(),
		Username: staff.Username,
		Email:    staff.Email,
		Role:     string(staff.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute)),
		},
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func VerifyToken(token, secretKey string) (*model.JWTPayload, error) {
	claims := &model.JWTClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return nil, err
	}

	payload := &model.JWTPayload{
		Id:       claims.Id,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     model.Role(claims.Role),
	}

	return payload, nil
}
