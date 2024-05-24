package controller

import (
	"beli-mang/service"
)

type Controller struct {
	svc service.Service
}

func NewController(svc service.Service) *Controller {
	return &Controller{
		svc: svc,
	}
}
