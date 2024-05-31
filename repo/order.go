package repo

import (
	"beli-mang/model"
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (model.Order, error)
	InsertCalculation(ctx context.Context, oc model.CalculatedEstimate) (model.CalculatedEstimate, error)
	GetCalculatedEstimateById(ctx context.Context, id uuid.UUID) (model.CalculatedEstimate, error)
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error
	GetUserOrders(ctx context.Context, params model.UserOrdersParams) ([]model.Order, error)
	GetNearbyMerchant(ctx context.Context, params model.GetMerchantParams, lat, long string) (listNearbyMerchant []model.GetNearbyMerchantData, meta model.MetaData, err error)
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order model.Order) (model.Order, error) {
	return order, nil
}

func (r *orderRepository) InsertCalculation(ctx context.Context, oc model.CalculatedEstimate) (model.CalculatedEstimate, error) {
	return model.CalculatedEstimate{}, nil
}
func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error {
	return nil
}

func (r *orderRepository) GetCalculatedEstimateById(ctx context.Context, id uuid.UUID) (model.CalculatedEstimate, error) {
	return model.CalculatedEstimate{}, nil
}

func (r *orderRepository) GetUserOrders(ctx context.Context, params model.UserOrdersParams) ([]model.Order, error) {
	return make([]model.Order, 0), nil
}

func (r *orderRepository) GetNearbyMerchant(ctx context.Context, params model.GetMerchantParams, lat, long string) (listNearbyMerchant []model.GetNearbyMerchantData, meta model.MetaData, err error) {
	floatLat, _ := strconv.ParseFloat(lat, 64)
	floatLong, _ := strconv.ParseFloat(long, 64)
	var getMerchantQuery = fmt.Sprintf(`SELECT *, earth_distance(ll_to_earth(latitude, longitude), ll_to_earth('%f', '%f')) AS distance FROM "merchant" WHERE true`, floatLat, floatLong)
	var total int = 0
	var metaData = model.MetaData{
		Offset: params.Offset,
		Limit:  params.Limit,
		Total:  0,
	}

	if params.Name != "" {
		name := "%" + params.Name + "%"
		getMerchantQuery += fmt.Sprintf(` AND "name" ILIKE '%s'`, name)
	}

	if params.MerchantId != "" {
		getMerchantQuery += fmt.Sprintf(` AND "id" = %s`, params.MerchantId)
	}

	if params.MerchantCategory != "" {
		getMerchantQuery += fmt.Sprintf(` AND "category" = %s`, params.MerchantCategory)
	}

	if params.CreatedAt != "" {
		if params.CreatedAt != "desc" && params.CreatedAt != "asc" {
			params.CreatedAt = "desc"
		}
		getMerchantQuery += fmt.Sprintf(` ORDER BY "createdAt" %s`, params.CreatedAt)
	} else {
		getMerchantQuery += ` ORDER BY "createdAt" DESC`
	}

	if params.Limit == 0 {
		params.Limit = 5 // default limit
	}

	getMerchantQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.Limit, params.Offset)
	getMerchantQuery += fmt.Sprintf(` ORDER BY distance ASC `)

	rows, err := r.db.QueryContext(ctx, getMerchantQuery)
	if err != nil {
		return nil, metaData, err
	}

	defer rows.Close()

	// Iterate over the rows and scan each row into a struct
	for rows.Next() {
		var merchant model.Merchant
		var items []model.Item
		var nearbyMerchant model.GetNearbyMerchantData
		if err := rows.Scan(&merchant.ID, &merchant.Name, &merchant.Category, &merchant.ImageURL, &merchant.Latitude, &merchant.Longitude, &merchant.CreatedAt); err != nil {
			return nil, metaData, err
		}
		var getItemById = `SELECT * FROM "merchantItem" WHERE "merchantId" = $1`
		rowsItem, err := r.db.QueryContext(ctx, getItemById, merchant.ID)
		if err != nil {
			return nil, metaData, err
		}

		for rowsItem.Next() {
			var item model.Item
			if err := rowsItem.Scan(&item.Id, &item.MerchantId, &item.Name, &item.Category, &item.ImageUrl, &item.Price, &item.CreatedAt); err != nil {
				return nil, metaData, err
			}

			items = append(items, item)
		}

		nearbyMerchant.Merchant = merchant
		nearbyMerchant.Items = items
		total += 1
		listNearbyMerchant = append(listNearbyMerchant, nearbyMerchant)
	}
	if err := rows.Err(); err != nil {
		return nil, metaData, err
	}

	metaData.Total = total

	return listNearbyMerchant, metaData, nil
}
