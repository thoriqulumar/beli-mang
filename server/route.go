package server

import (
	"beli-mang/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) RegisterRoute(cfg *config.Config) {
	mainRoute := s.app.Group("/v1")
	mainRoute.Any("/healthcheck", func(c echo.Context) error {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": http.StatusOK,
			"msg":  "ok",
		})
		return nil
	})
}
