package repo

import (
	"beli-mang/model"

	"github.com/jmoiron/sqlx"
)

type MerchantRepository interface {
	CreateMerchant(request model.Merchant) error
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
