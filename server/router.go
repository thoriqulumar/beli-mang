package server

import (
	"beli-mang/config"
	"beli-mang/controller"
	"beli-mang/repo"
	"beli-mang/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) RegisterRoute(cfg *config.Config) {
	mainRoute := s.app
	mainRoute.Any("/healthcheck", func(c echo.Context) error {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": http.StatusOK,
			"msg":  "ok",
		})
		return nil
	})

	registerImageRoute(mainRoute, cfg, s.logger)
	registerMerchantRoute(mainRoute, s.db, s.validator)
	registerStaffRoute(mainRoute, s.db, cfg, s.validator)
}

func registerImageRoute(e *echo.Echo, cfg *config.Config, logger *zap.Logger) {
	ctr := controller.NewImageController(service.NewImageService(cfg, logger))
	// auth := middleware.Authentication(cfg.JWTSecret)
	// e.POST("/image", auth(ctr.PostImage))
	// disable auth because it's not ready
	e.POST("/image", ctr.PostImage)
}

func registerMerchantRoute(e *echo.Echo, db *sqlx.DB, validate *validator.Validate) {
	ctr := controller.NewMerchantController(service.NewMerchantService(repo.NewMerchantRepository(db)), validate)

	e.POST("/admin/merchants", ctr.CreateMerchant)
	e.GET("/admin/merchants", ctr.GetMerchant)
	e.POST("/admin/merchants/:merchantId/items", ctr.CreateMerchantItem)
}

func registerPurchaseRoute(e *echo.Echo, db *sqlx.DB, validate *validator.Validate) {
	ctr := controller.NewPurchaseController(service.NewPurchaseService(repo.NewOrderRepository(db), repo.NewMerchantRepository(db)), validate)

	e.POST("/users/estimate", ctr.EstimateOrders)
	e.POST("/users/orders", ctr.ConfirmOrder)
	e.GET("/users/orders", ctr.GetUserOrders)
}

func registerStaffRoute(e *echo.Echo, db *sqlx.DB, cfg *config.Config, validate *validator.Validate) {
	ctr := controller.NewStaffController(service.NewStaffService(cfg, repo.NewStaffRepo(db)), validate)

	e.POST("/admin/register", ctr.RegisterStaffAdmin)
	e.POST("/admin/login", ctr.LoginStaffAdmin)

	e.POST("/users/register", ctr.RegisterStaffUser)
	e.POST("/users/login", ctr.LoginStaffUser)

}
