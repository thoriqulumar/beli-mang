package controller

import (
	"beli-mang/model"
	"beli-mang/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type MerchantController struct {
	svc      service.MerchantService
	validate *validator.Validate
}

func NewMerchantController(svc service.MerchantService, validate *validator.Validate) *MerchantController {
	_ = validate.RegisterValidation("custom_url", customURL)

	return &MerchantController{
		svc:      svc,
		validate: validate,
	}
}

func (ctr *MerchantController) CreateMerchant(ctx echo.Context) error {
	var createMerchantRequest model.CreateMerchantRequest
	if err := ctx.Bind(&createMerchantRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.CreateMerchantGeneralResponse{Message: "request doesn’t pass validation", Error: err.Error()})
	}

	merchantId, err := ctr.svc.CreateMerchant(createMerchantRequest)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.CreateMerchantGeneralResponse{Message: "Internal server error!", Error: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, model.CreateMerchantResponse{MerchantId: merchantId})
}
