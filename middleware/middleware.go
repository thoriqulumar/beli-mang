package middleware

import (
	"beli-mang/pkg/crypto"
	"beli-mang/pkg/customErr"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Authentication(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", -1)

			if token == "" {
				resErr := customErr.NewUnauthorizedError("Unauthorized")
				return c.JSON(resErr.StatusCode, resErr)
			}

			payload, err := crypto.VerifyToken(token, secret)
			if err != nil {
				resErr := customErr.NewUnauthorizedError("Unauthorized")
				if errors.Is(err, jwt.ErrTokenExpired) {
					resErr = customErr.NewUnauthorizedError("Token expired")
				}
				return c.JSON(resErr.StatusCode, resErr)
			}

			// Add user data to the request context
			c.Set("userData", payload)

			return next(c)
		}
	}
}
