package repo

import (
	"beli-mang/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"strconv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	Create(ctx context.Context, tx *sqlx.Tx, order model.Order) (model.Order, error)
	InsertCalculation(ctx context.Context, tx *sqlx.Tx, oc model.CalculatedEstimate) (model.CalculatedEstimate, error)
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

func (r *orderRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *orderRepository) Create(ctx context.Context, tx *sqlx.Tx, order model.Order) (model.Order, error) {
	var createOrderQuery = `INSERT INTO "order"(
		"orderId",
		"orderStatus",
		detail,
		"merchantIds",
		"joinedMerchantName",
		"merchantCategories",
		"joinedItemsName",
		"userId",
		"userLatitude",
		"userLongitude",
		"createdAt"
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`

	_, err := tx.ExecContext(ctx, createOrderQuery,
		order.OrderID,
		order.OrderStatus,
		order.DetailRaw,
		pq.Array(order.MerchantIDs),
		order.JoinedMerchantName,
		pq.Array(order.MerchantCategories),
		order.JoinedItemsName,
		order.UserID,
		order.UserLatitude,
		order.UserLongitude,
		order.CreatedAt,
	)

	return order, err
}

func (r *orderRepository) InsertCalculation(ctx context.Context, tx *sqlx.Tx, oc model.CalculatedEstimate) (model.CalculatedEstimate, error) {
	var insertCalculationQuery = `INSERT INTO "calculatedEstimate"(
		"calculatedEstimateId",
		"totalPrice",
		"estimatedDeliveryTimeInMinutes",
		"orderId",
		createdAt)
	VALUES($1, $2, $3, $4, $5);
	`
	_, err := tx.ExecContext(ctx, insertCalculationQuery,
		oc.CalculatedEstimateId,
		oc.TotalPrice,
		oc.EstimatedDeliveryTimeInMinutes,
		oc.OrderId,
		oc.CreatedAt)
	return oc, err
}
func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error {
	var updateOrderStatus = `UPDATE "order" SET "orderStatus"=$2 WHERE "orderId"=$1`
	_, err := r.db.ExecContext(ctx, updateOrderStatus, orderID, status)
	return err
}

func (r *orderRepository) GetCalculatedEstimateById(ctx context.Context, id uuid.UUID) (model.CalculatedEstimate, error) {
	var getCalculatedEstimateByIdQuery = `SELECT * FROM "calculatedEstimate" WHERE "calculatedEstimateId" = $1`
	var result model.CalculatedEstimate
	err := r.db.QueryRowxContext(ctx, getCalculatedEstimateByIdQuery, id).StructScan(result)
	return result, err
}

func (r *orderRepository) GetUserOrders(ctx context.Context, params model.UserOrdersParams) ([]model.Order, error) {
	var listOrder []model.Order
	var getUserOrdersQuery = `SELECT * FROM "order" WHERE "userId" = $1 AND "orderStatus" = $2 AND TRUE`
	if params.MerchantId != nil {
		getUserOrdersQuery += `AND "merchantId"=` + params.MerchantId.String()
	}
	if params.Name != nil {
		// query with searchable index
		nameStmt := fmt.Sprintf(`to_tsvector('english', "joinedMerchantName") @@ plainto_tsquery('english', %s)
   OR to_tsvector('english', "joinedItemsName") @@ plainto_tsquery('english', %s)`, *params.Name)
		getUserOrdersQuery += `AND (` + nameStmt + `)`
	}

	if params.MerchantCategory != nil {
		getUserOrdersQuery += fmt.Sprintf(`AND "merchantCategories"=ANY(%s)`, *params.MerchantCategory)
	}

	getUserOrdersQuery += fmt.Sprintf(`ORDER BY "createdAt" DESC LIMIT %d OFFSET %d`, params.Limit, params.Offset)
	rows, err := r.db.QueryxContext(ctx, getUserOrdersQuery, params.UserID, params.Status)
	if err != nil {
		return listOrder, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		var order model.Order
		var detail model.OrderDetail
		if err := rows.StructScan(&order); err != nil {
			return listOrder, err
		}
		_ = json.Unmarshal(order.DetailRaw, &detail)
		order.Detail = detail
		listOrder = append(listOrder, order)
	}
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
		getMerchantQuery += fmt.Sprintf(` AND "id" = '%s'`, params.MerchantId)
	}

	if params.MerchantCategory != "" {
		getMerchantQuery += fmt.Sprintf(` AND "category" = '%s'`, params.MerchantCategory)
	}

	orderClause := ""
	if params.CreatedAt != "" {
		if params.CreatedAt != "desc" && params.CreatedAt != "asc" {
			params.CreatedAt = "desc"
		}
		orderClause = fmt.Sprintf(` ORDER BY "createdAt" %s`, params.CreatedAt)
	} else {
		orderClause = ` ORDER BY "createdAt" DESC`
	}

	if orderClause != "" {
		orderClause = fmt.Sprintf(` ORDER BY distance ASC `)
	} else {
		orderClause += `,distance ASC`
	}

	if params.Limit == 0 {
		params.Limit = 5 // default limit
	}
	getMerchantQuery += orderClause
	getMerchantQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.Limit, params.Offset)

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
		var distance float64
		if err := rows.Scan(&merchant.ID, &merchant.Name, &merchant.Category, &merchant.ImageURL, &merchant.Location.Lat, &merchant.Location.Long, &merchant.CreatedAt, &distance); err != nil {
			return nil, metaData, err
		}
		var getItemById = `SELECT * FROM "merchantItem" WHERE "merchantId" = $1`
		rowsItem, err := r.db.QueryContext(ctx, getItemById, merchant.ID)
		if err != nil {
			return nil, metaData, err
		}

		defer rowsItem.Close()

		items = []model.Item{}
		for rowsItem.Next() {
			var item model.Item
			if err := rowsItem.Scan(&item.Id, &item.MerchantId, &item.Name, &item.Category, &item.ImageUrl, &item.Price, &item.CreatedAt); err != nil {
				return nil, metaData, err
			}

			items = append(items, item)
		}

		nearbyMerchant.Merchant = merchant
		nearbyMerchant.Items = items
		nearbyMerchant.Distance = fmt.Sprintf("%.2f m", distance)
		total += 1
		listNearbyMerchant = append(listNearbyMerchant, nearbyMerchant)
	}
	if err := rows.Err(); err != nil {
		return nil, metaData, err
	}

	metaData.Total = total
	metaData.Offset = params.Offset
	metaData.Limit = params.Limit

	return listNearbyMerchant, metaData, nil
}
