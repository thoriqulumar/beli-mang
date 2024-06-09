package model

import (
	"github.com/google/uuid"
	"time"
)

type ItemCategory string

type Item struct {
	Id         uuid.UUID    `json:"id" db:"id"`
	MerchantId uuid.UUID    `json:"merchantId" db:"merchantId"`
	Name       string       `json:"name" db:"name"`
	Category   ItemCategory `json:"category" db:"category"`
	ImageUrl   string       `json:"imageUrl" db:"imageUrl"`
	Price      int          `json:"price" db:"price"`
	CreatedAt  time.Time    `json:"createdAt" db:"createdAt"`
}

func (i Item) ItemToBoughtItem() BoughtItem {
	return BoughtItem{
		ItemId:          i.Id,
		Name:            i.Name,
		ProductCategory: string(i.Category),
		Price:           i.Price,
		ImageUrl:        i.ImageUrl,
		CreatedAt:       i.CreatedAt,
		Quantity:        0,
	}
}
