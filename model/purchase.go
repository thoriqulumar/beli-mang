package model

import (
	"github.com/google/uuid"
	"time"
)

type EstimateOrdersRequest struct {
	UserId       uuid.UUID      `json:"userId"`
	UserLocation UserLocation   `json:"userLocation"`
	Orders       []OrderRequest `json:"orders"`
}

type EstimateOrdersResponse struct {
	TotalPrice                     int       `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int       `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId           uuid.UUID `json:"calculatedEstimateId"`
}

type CalculatedEstimate struct {
	TotalPrice                     int       `json:"totalPrice" db:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int       `json:"estimatedDeliveryTimeInMinutes" db:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId           uuid.UUID `json:"calculatedEstimateId" db:"calculatedEstimateId"`
	OrderId                        uuid.UUID `json:"orderId" db:"orderId"`
	CreatedAt                      time.Time `json:"createdAt" db:"createdAt"`
}

type UserLocation struct {
	Lat  float64 `json:"lat"`  //Latitude
	Long float64 `json:"long"` //Longitude
}

type ConfirmOrderRequest struct {
	CalculatedEstimateId uuid.UUID `json:"calculatedEstimateId"`
}

type ConfirmOrderResponse struct {
	OrderId uuid.UUID `json:"orderId"`
}

type GetUserOrdersRequest struct {
	MerchantId *string
	Limit      *int
	Offset     *int
	Name       *string
	Category   *MerchantCategory
}

// Point represents a geographical point with latitude and longitude
type Point struct {
	Lat float64 // Latitude
	Lon float64 // Longitude
}
