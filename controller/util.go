package controller

import (
	"beli-mang/model"
	"github.com/labstack/echo/v4"
)

func GetUserFromContext(c echo.Context) *model.JWTPayload {
	jwtPayload := c.Get("userData")
	return jwtPayload.(*model.JWTPayload)
}
