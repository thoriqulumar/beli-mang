package controller

import (
	"beli-mang/model"
	"beli-mang/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PurchaseController struct {
	svc      service.PurchaseService
	validate *validator.Validate
}

/*
	e.POST("/users/estimate", ctr.EstimateOrders)
	e.POST("/users/orders", ctr.ConfirmOrder)
	e.GET("/users/orders", ctr.GetUserOrders)
*/

func NewPurchaseController(svc service.PurchaseService, validate *validator.Validate) *PurchaseController {
	_ = validate.RegisterValidation("custom_url", customURL)

	return &PurchaseController{
		svc:      svc,
		validate: validate,
	}
}

func (ctr *PurchaseController) EstimateOrders(ctx echo.Context) error {
	var etaOrderRequest model.EstimateOrdersRequest
	if err := ctx.Bind(&etaOrderRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: "request doesn’t pass validation", Error: err.Error()})
	}

	data, err := ctr.svc.EstimateOrders(ctx.Request().Context(), etaOrderRequest)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.GeneralResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, data)
}

func (ctr *PurchaseController) ConfirmOrder(ctx echo.Context) error {
	var payload model.ConfirmOrderRequest
	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: "request doesn’t pass validation", Error: err.Error()})
	}
	data, err := ctr.svc.ConfirmOrder(ctx.Request().Context(), payload)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.GeneralResponse{
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, data)
}

func (ctr *PurchaseController) GetUserOrders(ctx echo.Context) error {
	var params model.UserOrdersParams

	data, err := ctr.svc.GetUserOrders(ctx.Request().Context(), params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.GeneralResponse{
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, data)
}
