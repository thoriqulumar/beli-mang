package controller

import (
	"beli-mang/model"
	"beli-mang/service"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

	if err := ctr.validate.Struct(createMerchantRequest); err != nil {
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
		Meta:    meta,
	})
}

func (ctr *MerchantController) CreateMerchantItem(ctx echo.Context) error {
	merchantID := ctx.Param("merchantId")
	var createMerchantItemRequest model.CreateMerchantItemRequest
	if err := ctx.Bind(&createMerchantItemRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.CreateMerchantGeneralResponse{Message: "request doesn’t pass validation", Error: err.Error()})
	}

	if err := ctr.validate.Struct(createMerchantItemRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.CreateMerchantGeneralResponse{Message: "request doesn’t pass validation", Error: err.Error()})
	}

	merchantUUID, err := uuid.Parse(merchantID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.CreateMerchantGeneralResponse{Message: "Internal server error!", Error: err.Error()})
	}

	itemId, err := ctr.svc.CreateMerchantItem(ctx.Request().Context(), createMerchantItemRequest, merchantUUID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.CreateMerchantGeneralResponse{Message: "Internal server error!", Error: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, model.CreateMerchantItemResponse{ItemId: itemId})
}

func (ctr *MerchantController) GetMerchantItem(ctx echo.Context) error {
	merchantID := ctx.Param("merchantId")
	value, err := ctx.FormParams()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	merchantUUID, err := uuid.Parse(merchantID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.CreateMerchantGeneralResponse{Message: "Internal server error!", Error: err.Error()})
	}

	// query to service
	data, meta, err := ctr.svc.GetMerchantItem(ctx.Request().Context(), merchantUUID, parseGetMerchantItemParams(value))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, model.MerchantGeneralResponse{
		Message: "success",
		Data:    data,
		Meta:    meta,
	})
}

func parseGetMerchantItemParams(params url.Values) model.GetMerchantItemParams {
	var result model.GetMerchantItemParams

	for key, values := range params {
		switch key {
		case "itemId":
			result.ItemId = values[0]
		case "name":
			result.Name = values[0]
		case "productCategory":
			result.ProductCategory = values[0]
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

func parseGetMerchantParams(params url.Values) model.GetMerchantParams {
	var result model.GetMerchantParams

	for key, values := range params {
		switch key {
		case "merchantId":
			result.MerchantId = values[0]
		case "name":
			result.Name = values[0]
		case "merchantCategory":
			result.MerchantCategory = values[0]
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
