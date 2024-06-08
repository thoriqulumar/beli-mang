package model

import (
	"time"

	"github.com/google/uuid"
)

type MerchantCategory string

// enum of merchant category, reduce miss typing
const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

type Merchant struct {
	ID        uuid.UUID        `json:"merchantId" db:"id"`
	Name      string           `json:"name" db:"name"`
	Category  MerchantCategory `json:"merchantCategory" db:"category"`
	ImageURL  string           `json:"imageUrl" db:"imageUrl"`
	Location  Location         `json:"location" db:"-"`
	CreatedAt time.Time        `json:"createdAt" db:"createdAt"`
}

type CreateMerchantRequest struct {
	Name     string   `json:"name" validate:"required,min=2,max=30"`
	Category string   `json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
	ImageURL string   `json:"imageUrl" validate:"required,custom_url"`
	Location Location `json:"location" validate:"required"`
}

type Location struct {
	Lat  float64 `json:"lat" validate:"required" db:"latitude"`
	Long float64 `json:"long" validate:"required" db:"longitude"`
}

type CreateMerchantResponse struct {
	MerchantId string `json:"merchantId"`
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
	ID         uuid.UUID `json:"itemId" db:"id"`
	MerchantId uuid.UUID `json:"merchantId" db:"merchantId"`
	Name       string    `json:"name" db:"name"`
	Category   string    `json:"productCategory" db:"category"`
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

type GetMerchantItemParams struct {
	ItemId          string
	Name            string
	ProductCategory string
	Limit           int
	Offset          int
	CreatedAt       string
}

type GetNearbyMerchantData struct {
	Merchant Merchant `json:"merchant"`
	Items    []Item   `json:"items"`
	Distance string   `json:"distance"`
}
