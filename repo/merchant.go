package repo

import (
	"beli-mang/model"
	"context"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type MerchantRepository interface {
	CreateMerchant(request model.Merchant) error
	GetMerchantMapByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]model.Merchant, error)
	GetMerchantItemMapByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]model.Item, error)
}

type merchantRepository struct {
	db *sqlx.DB
}

func NewMerchantRepository(db *sqlx.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

var (
	createMerchantQuery = `
	INSERT INTO merchant (id, name, category, "imageUrl", latitude, longitude, "createdAt")
	VALUES ($1, $2, $3, $4, $5, $6, NOW())
	RETURNING id;
`
)

func (r *merchantRepository) CreateMerchant(request model.Merchant) error {
	return r.db.QueryRowx(createMerchantQuery, request.ID, request.Name, request.Category, request.ImageURL, request.Latitude, request.Longitude).Scan(&request.ID)
}

func (r *merchantRepository) GetMerchantMapByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]model.Merchant, error) {
	var merchants map[uuid.UUID]model.Merchant
	var getMerchantQuery = `SELECT * FROM merchant WHERE id IN ($1)`
	rows, err := r.db.QueryxContext(ctx, getMerchantQuery, pq.Array(ids))
	if err != nil {
		return merchants, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var merchant model.Merchant
		if err := rows.StructScan(&merchant); err != nil {
			return merchants, err
		}
		merchants[merchant.ID] = merchant
	}
	if err := rows.Err(); err != nil {
		return merchants, err
	}
	return merchants, nil
}

func (r *merchantRepository) GetMerchantItemMapByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]model.Item, error) {
	mapItems := make(map[uuid.UUID]model.Item)
	var getItemQuery = `SELECT * FROM "merchantItem" WHERE id IN ($1)`
	rows, err := r.db.QueryxContext(ctx, getItemQuery, pq.Array(ids))
	if err != nil {
		return mapItems, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		var merchantItem model.Item
		if err := rows.StructScan(&merchantItem); err != nil {
			return mapItems, err
		}
		mapItems[merchantItem.Id] = merchantItem
	}

	return mapItems, nil
}
