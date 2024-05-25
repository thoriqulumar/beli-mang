package repo

import (
	"beli-mang/model"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MerchantRepository interface {
	CreateMerchant(request model.Merchant) error
	GetMerchant(ctx context.Context, params model.GetMerchantParams) (patients []model.Merchant, meta model.MetaData, err error)
	GetMerchantById(ctx context.Context, merchantId string) (merchant model.Merchant, err error)
	CreateMerchantItem(request model.MerchantItem) error
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

var (
	getMerchantByIdQuery = `
	SELECT * FROM "merchant" WHERE id = $1;
`
)

func (r *merchantRepository) GetMerchantById(ctx context.Context, merchantId string) (merchant model.Merchant, err error) {
	err = r.db.QueryRowxContext(ctx, getMerchantByIdQuery, merchantId).StructScan(&merchant)

	if err != nil {
		return
	}

	return merchant, nil
}

func (r *merchantRepository) GetMerchant(ctx context.Context, params model.GetMerchantParams) (patients []model.Merchant, meta model.MetaData, err error) {
	var listMerchant []model.Merchant
	var getMerchantQuery = `SELECT * FROM "merchant" WHERE true`
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

	rows, err := r.db.QueryContext(ctx, getMerchantQuery)
	if err != nil {
		return nil, metaData, err
	}

	defer rows.Close()

	// Iterate over the rows and scan each row into a struct
	for rows.Next() {
		var merchant model.Merchant
		if err := rows.Scan(&merchant.ID, &merchant.Name, &merchant.Category, &merchant.ImageURL, &merchant.Latitude, &merchant.Longitude, &merchant.CreatedAt); err != nil {
			return nil, metaData, err
		}
		total += 1
		listMerchant = append(listMerchant, merchant)
	}
	if err := rows.Err(); err != nil {
		return nil, metaData, err
	}

	metaData.Total = total

	return listMerchant, metaData, nil
}


var (
	createMerchantItemQuery = `
	INSERT INTO "merchantItem" (id, "merchantId", name, "category", "imageUrl", price, "createdAt")
	VALUES ($1, $2, $3, $4, $5, $6, NOW())
	RETURNING id;
`
)

func (r *merchantRepository) CreateMerchantItem(request model.MerchantItem) error {
	return r.db.QueryRowx(createMerchantItemQuery, request.ID, request.MerchantId, request.Name, request.Category, request.ImageURL, request.Price).Scan(&request.ID)
}