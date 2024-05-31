package server

import (
	"beli-mang/config"
	"beli-mang/controller"
	"beli-mang/middleware"
	"beli-mang/model"
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
	registerMerchantRoute(mainRoute, s.db, cfg, s.validator)
	registerStaffRoute(mainRoute, s.db, cfg, s.validator)
	registerPurchaseRoute(mainRoute, s.db, cfg, s.validator, s.logger)
}

func registerImageRoute(e *echo.Echo, cfg *config.Config, logger *zap.Logger) {
	ctr := controller.NewImageController(service.NewImageService(cfg, logger))
	auth := middleware.Authentication(cfg.JWTSecret, model.RoleAdmin)
	// e.POST("/image", auth(ctr.PostImage))
	// disable auth because it's not ready
	e.POST("/image", auth(ctr.PostImage))
}

func registerMerchantRoute(e *echo.Echo, db *sqlx.DB, cfg *config.Config, validate *validator.Validate) {
	ctr := controller.NewMerchantController(service.NewMerchantService(repo.NewMerchantRepository(db)), validate)

	auth := middleware.Authentication(cfg.JWTSecret, model.RoleAdmin)
	e.POST("/admin/merchants", auth(ctr.CreateMerchant))
	e.GET("/admin/merchants", auth(ctr.GetMerchant))
	e.POST("/admin/merchants/:merchantId/items", auth(ctr.CreateMerchantItem))
	e.GET("/admin/merchants/:merchantId/items", auth(ctr.GetMerchantItem))
}

func registerPurchaseRoute(e *echo.Echo, db *sqlx.DB, cfg *config.Config, validate *validator.Validate, logger *zap.Logger) {
	ctr := controller.NewPurchaseController(service.NewPurchaseService(repo.NewOrderRepository(db), repo.NewMerchantRepository(db), logger), validate)

	auth := middleware.Authentication(cfg.JWTSecret, model.RoleUser)
	e.POST("/users/estimate", auth(ctr.EstimateOrders))
	e.POST("/users/orders", auth(ctr.ConfirmOrder))
	e.GET("/users/orders", auth(ctr.GetUserOrders))
}

func registerStaffRoute(e *echo.Echo, db *sqlx.DB, cfg *config.Config, validate *validator.Validate) {
	ctr := controller.NewStaffController(service.NewStaffService(cfg, repo.NewStaffRepo(db)), validate)

	e.POST("/admin/register", ctr.RegisterStaffAdmin)
	e.POST("/admin/login", ctr.LoginStaffAdmin)

	e.POST("/users/register", ctr.RegisterStaffUser)
	e.POST("/users/login", ctr.LoginStaffUser)

}
