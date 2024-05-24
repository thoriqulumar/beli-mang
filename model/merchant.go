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
