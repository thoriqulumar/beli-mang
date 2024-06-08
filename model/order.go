package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type OrderRequest struct {
	MerchantId      string             `json:"merchantId" validate:"required"`
	IsStartingPoint bool               `json:"isStartingPoint"`
	Items           []OrderRequestItem `json:"items" validate:"required,dive"`
}

type OrderRequestItem struct {
	ItemId   string `json:"itemId" validate:"required"`
	Quantity int    `json:"quantity"`
}

type GetUserOrdersResponse []UserOrderData

type UserOrderData struct {
	OrderId uuid.UUID   `json:"orderId"`
	Orders  []OrderData `json:"orders"`
}

type OrderData struct {
	Merchant        Merchant     `json:"merchant"`
	IsStartingPoint bool         `json:"isStartingPoint,omitempty"`
	Items           []BoughtItem `json:"items"`
}

// Item ...
// TODO: remove this if items endpoint done. the struct declaration should be on item.go file
type BoughtItem struct {
	ItemId          uuid.UUID `json:"itemId"`
	Name            string    `json:"name"`
	ProductCategory string    `json:"productCategory"`
	Price           int       `json:"price"`
	ImageUrl        string    `json:"imageUrl"`
	CreatedAt       string    `json:"createdAt"`
	Quantity        int       `json:"quantity"`
}

// OrderStatus enum type
type OrderStatus string

const (
	OrderStatusDraft   OrderStatus = "DRAFT"
	OrderStatusCreated OrderStatus = "CREATED"
)

// Order struct
type Order struct {
	OrderID            uuid.UUID          `json:"orderId" db:"orderId"`
	OrderStatus        OrderStatus        `json:"orderStatus" db:"orderStatus"`
	DetailRaw          json.RawMessage    `json:"-" db:"detail"`
	Detail             OrderDetail        `json:"detail" db:"-"`
	MerchantIDs        []uuid.UUID        `json:"merchantIds" db:"merchantIds"`
	JoinedMerchantName string             `json:"joinedMerchantName" db:"joinedMerchantName"`
	MerchantCategories []MerchantCategory `json:"merchantCategories" db:"merchantCategories"`
	JoinedItemsName    string             `json:"joinedItemsName" db:"joinedItemsName"`
	UserID             uuid.UUID          `json:"userId" db:"userId"`
	UserLatitude       float64            `json:"userLatitude" db:"userLatitude"`
	UserLongitude      float64            `json:"userLongitude" db:"userLongitude"`
	CreatedAt          time.Time          `json:"createdAt" db:"createdAt"`
}

func (o Order) ToUserOrderData() UserOrderData {
	return UserOrderData{
		OrderId: o.OrderID,
		Orders:  o.Detail,
	}
}

type OrderDetail []OrderData

type UserOrdersParams struct {
	UserID           uuid.UUID         `json:"userId"`
	Status           OrderStatus       `json:"status"`
	MerchantId       *uuid.UUID        `json:"merchantId"`
	Limit            *int              `json:"limit"`
	Offset           *int              `json:"offset"`
	Name             *string           `json:"name"`
	MerchantCategory *MerchantCategory `json:"merchantCategory"`
}
