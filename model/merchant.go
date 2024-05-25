package model

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Category  string    `json:"category" db:"category"`
	ImageURL  string    `json:"imageUrl" db:"imageUrl"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}

type CreateMerchantRequest struct {
	Name      string  `json:"name" validate:"required,min=2,max=30"`
	Category  string  `json:"category" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
	ImageURL  string  `json:"imageUrl" validate:"required,custom_url"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type CreateMerchantResponse struct {
	MerchantId string
}

type CreateMerchantGeneralResponse struct {
	Message string
	Error   string
}

type MetaData struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}
type MerchantGeneralResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    MetaData    `json:"meta,omitempty"`
}

type GetMerchantParams struct {
	MerchantId       string
	Name             string
	MerchantCategory string
	Limit            int
	Offset           int
	CreatedAt        string
}

type MerchantItem struct {
	ID         uuid.UUID `json:"id" db:"id"`
	MerchantId string    `json:"merchantId" db:"merchantId"`
	Name       string    `json:"name" db:"name"`
	Category   string    `json:"category" db:"category"`
	ImageURL   string    `json:"imageUrl" db:"imageUrl"`
	Price      int       `json:"price" db:"price"`
	CreatedAt  time.Time `json:"createdAt" db:"createdAt"`
}

type CreateMerchantItemRequest struct {
	Name            string `json:"name" validate:"required,min=2,max=30"`
	ProductCategory string `json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
	ImageURL        string `json:"imageUrl" validate:"required,custom_url"`
	Price           int    `json:"price" validate:"required"`
}

type CreateMerchantItemResponse struct {
	ItemId string `json:"itemId"`
}
