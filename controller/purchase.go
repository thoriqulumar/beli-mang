package controller

import (
	"beli-mang/model"
	cerr "beli-mang/pkg/customErr"
	"beli-mang/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: "invalid format payload", Error: err.Error()})
	}

	if err := ctr.validate.Struct(etaOrderRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: "request doesn’t pass validation", Error: err.Error()})
	}

	user := GetUserFromContext(ctx)
	etaOrderRequest.UserId, _ = uuid.Parse(user.Id)
	data, err := ctr.svc.EstimateOrders(ctx.Request().Context(), etaOrderRequest)
	if err != nil {
		errCode := cerr.GetCode(err)
		if errCode == 0 {
			errCode = http.StatusInternalServerError
		}
		return ctx.JSON(errCode, model.GeneralResponse{
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

	if err := ctr.validate.Struct(payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: "request doesn’t pass validation", Error: err.Error()})
	}
	data, err := ctr.svc.ConfirmOrder(ctx.Request().Context(), payload)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.GeneralResponse{
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusCreated, data)
}

func (ctr *PurchaseController) GetUserOrders(ctx echo.Context) error {
	var params model.UserOrdersParams
	value, err := ctx.FormParams()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	params = parseUserOrdersParams(value)
	user := GetUserFromContext(ctx)
	params.UserID, _ = uuid.Parse(user.Id)

	data, err := ctr.svc.GetUserOrders(ctx.Request().Context(), params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.GeneralResponse{
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, data)
}

func (ctr *PurchaseController) GetMerchantNearby(ctx echo.Context) error {
	latlong := ctx.Param("latlong")
	temp := strings.Split(latlong, ",")
	if len(temp) != 2 {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "lat/long not valid"})
	}
	lat := temp[0]
	long := temp[1]
	err := ValidateLatLong(lat, long)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "lat/long not valid"})
	}

	value, err := ctx.FormParams()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	// query to service
	data, meta, err := ctr.svc.GetNearbyMerchant(ctx.Request().Context(), parseGetMerchantParams(value), lat, long)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, model.MerchantGeneralResponse{
		Message: "success",
		Data:    data,
		Meta:    meta,
	})
}

func parseUserOrdersParams(params url.Values) model.UserOrdersParams {
	var result model.UserOrdersParams
	for key, values := range params {
		switch key {
		case "name":
			n := values[0]
			result.Name = &n
		case "status":
			s := model.OrderStatus(values[0])
			result.Status = s
		case "merchantCategory":
			v := model.MerchantCategory(values[0])
			result.MerchantCategory = &v
		case "limit":
			limit, err := strconv.Atoi(values[0])
			if err == nil {
				result.Limit = &limit
			}
		case "offset":
			offset, err := strconv.Atoi(values[0])
			if err == nil {
				result.Offset = &offset
			}
		}
	}

	// set default
	if result.Limit == nil {
		defaultLimit := 5
		result.Limit = &defaultLimit
	}

	return result
}
