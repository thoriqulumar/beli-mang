package server

import (
	"beli-mang/config"
	"beli-mang/controller"
	"beli-mang/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) RegisterRoute(cfg *config.Config) {
	mainRoute := s.app

	registerImageRoute(mainRoute, cfg, s.logger)
}

func registerImageRoute(e *echo.Echo, cfg *config.Config, logger *zap.Logger) {
	ctr := controller.NewImageController(service.NewImageService(cfg, logger))
	// auth := middleware.Authentication(cfg.JWTSecret)
	// e.POST("/image", auth(ctr.PostImage))
	// disable auth because it's not ready
	e.POST("/image", ctr.PostImage)
}
