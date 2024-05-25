package repo

import (
	"beli-mang/model"
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (model.Order, error)
	InsertCalculation(ctx context.Context, oc model.CalculatedEstimate) (model.CalculatedEstimate, error)
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
