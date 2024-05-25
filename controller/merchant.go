package controller

import (
	"beli-mang/model"
	"beli-mang/service"
	"net/http"
	"net/url"
	"strconv"

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

func (ctr *MerchantController) GetMerchant(ctx echo.Context) error {
	value, err := ctx.FormParams()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	// query to service
	data, meta, err := ctr.svc.GetMerchant(ctx.Request().Context(), parseGetMerchantParams(value))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, model.MerchantGeneralResponse{
		Message: "success",
		Data:    data,
		Meta: &meta,
	})
}


func parseGetMerchantParams(params url.Values) model.GetMerchantParams {
	var result model.GetMerchantParams

	for key, values := range params {
		switch key {
		case "merchantId":
			result.MerchantId = values[0]
		case "name":
			result.Name = values[0]
		case "limit":
			limit, err := strconv.Atoi(values[0])
			if err == nil {
				result.Limit = limit
			}
		case "offset":
			offset, err := strconv.Atoi(values[0])
			if err == nil {
				result.Offset = offset
			}
		case "createdAt":
			result.CreatedAt = values[0]
		}
	}

	return result
}
