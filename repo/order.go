package repo

import (
	"beli-mang/model"
	"context"
	"encoding/json"
	"fmt"
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
		order.MerchantIDs,
		order.JoinedMerchantName,
		order.MerchantCategories,
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
