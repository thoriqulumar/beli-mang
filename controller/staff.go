package controller

import (
	"beli-mang/model"
	cerr "beli-mang/pkg/customErr"
	"beli-mang/service"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

type StaffController struct {
	svc      service.StaffService
	validate *validator.Validate
}

func NewStaffController(svc service.StaffService, validate *validator.Validate) *StaffController {
	_ = validate.RegisterValidation("email", isEmailValid)
	return &StaffController{
		svc:      svc,
		validate: validate,
	}
}

func (c *StaffController) RegisterStaffAdmin(ctx echo.Context) error {
	var newStaffReq model.RegisterStaffRequest
	if err := ctx.Bind(&newStaffReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	err := c.validate.Struct(newStaffReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}

	newStaff := model.Staff{
		Username: newStaffReq.Username,
		Password: newStaffReq.Password,
		Email:    newStaffReq.Email,
		Role:     model.RoleAdmin,
	}
	serviceRes, err := c.svc.RegisterAdmin(ctx.Request().Context(), newStaff)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})

	}
	registerStaffResponse := model.StaffWithToken{
		AccessToken: serviceRes.AccessToken,
	}
	return ctx.JSON(http.StatusCreated, registerStaffResponse)
}

func (c *StaffController) LoginStaffAdmin(ctx echo.Context) error {
	var staffReq model.LoginStaffRequest
	if err := ctx.Bind(&staffReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}
	err := c.validate.Struct(staffReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}
	staff := model.Staff{
		Username: staffReq.Username,
		Password: staffReq.Password,
	}
	serviceRes, err := c.svc.LoginAdmin(ctx.Request().Context(), staff)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})
	}
	loginStaffResponse := model.StaffWithToken{
		AccessToken: serviceRes.AccessToken,
	}
	return ctx.JSON(http.StatusOK, loginStaffResponse)
}

func (c *StaffController) RegisterStaffUser(ctx echo.Context) error {
	var newStaffReq model.RegisterStaffRequest
	if err := ctx.Bind(&newStaffReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	err := c.validate.Struct(newStaffReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}

	newStaff := model.Staff{
		Username: newStaffReq.Username,
		Password: newStaffReq.Password,
		Email:    newStaffReq.Email,
		Role:     model.RoleAdmin,
	}
	serviceRes, err := c.svc.RegisterUser(ctx.Request().Context(), newStaff)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})

	}
	registerStaffResponse := model.StaffWithToken{
		AccessToken: serviceRes.AccessToken,
	}
	return ctx.JSON(http.StatusCreated, registerStaffResponse)
}

func (c *StaffController) LoginStaffUser(ctx echo.Context) error {
	var staffReq model.LoginStaffRequest
	if err := ctx.Bind(&staffReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}
	err := c.validate.Struct(staffReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}
	staff := model.Staff{
		Username: staffReq.Username,
		Password: staffReq.Password,
	}
	serviceRes, err := c.svc.LoginAdmin(ctx.Request().Context(), staff)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})
	}
	loginStaffResponse := model.StaffWithToken{
		AccessToken: serviceRes.AccessToken,
	}
	return ctx.JSON(http.StatusOK, loginStaffResponse)
}
