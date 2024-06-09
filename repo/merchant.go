package repo

import (
	"beli-mang/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MerchantRepository interface {
	CreateMerchant(request model.Merchant) error
	GetMerchantMapByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]model.Merchant, error)
	GetMerchantItemMapByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]model.Item, error)
	GetMerchant(ctx context.Context, params model.GetMerchantParams) (patients []model.Merchant, meta model.MetaData, err error)
	GetMerchantById(ctx context.Context, merchantId uuid.UUID) (merchant model.Merchant, err error)
	CreateMerchantItem(request model.MerchantItem) error
	GetMerchantItem(ctx context.Context, params model.GetMerchantItemParams) (patients []model.MerchantItem, meta model.MetaData, err error)
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
	return r.db.QueryRowx(createMerchantQuery, request.ID, request.Name, request.Category, request.ImageURL, request.Location.Lat, request.Location.Long).Scan(&request.ID)
}

func (r *merchantRepository) GetMerchantMapByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]model.Merchant, error) {
	merchants := make(map[uuid.UUID]model.Merchant)
	// Create a slice of placeholders and convert UUIDs to their string representation
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id.String()
	}

	// Join the placeholders with commas to form the IN clause
	pStr := fmt.Sprintf("IN (%s)", strings.Join(placeholders, ", "))
	getMerchantQuery := `SELECT * FROM merchant WHERE id ` + pStr

	rows, err := r.db.QueryxContext(ctx, getMerchantQuery, args...)
	if err != nil {
		return merchants, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var merchant model.Merchant
		if err := rows.Scan(&merchant.ID, &merchant.Name, &merchant.Category, &merchant.ImageURL, &merchant.Location.Lat, &merchant.Location.Long, &merchant.CreatedAt); err != nil {
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
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id.String() // Convert UUID to string
	}

	// Join the placeholders with commas to form the IN clause
	pStr := fmt.Sprintf("IN (%s)", strings.Join(placeholders, ", "))
	var getItemQuery = `SELECT * FROM "merchantItem" WHERE id ` + pStr
	rows, err := r.db.QueryxContext(ctx, getItemQuery, args...)
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

var (
	getMerchantByIdQuery = `
	SELECT "id",
     "name",
     "category",
     "imageUrl",
     "latitude",
     "longitude",
     "createdAt" FROM "merchant" WHERE id = $1;
`
)

func (r *merchantRepository) GetMerchantById(ctx context.Context, merchantId uuid.UUID) (merchant model.Merchant, err error) {
	err = r.db.QueryRowxContext(ctx, getMerchantByIdQuery, merchantId).
		Scan(&merchant.ID, &merchant.Name, &merchant.Category, &merchant.ImageURL, &merchant.Location.Lat, &merchant.Location.Long, &merchant.CreatedAt)
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
		getMerchantQuery += fmt.Sprintf(` AND "id" = '%s'`, params.MerchantId)
	}

	if params.MerchantCategory != "" {
		getMerchantQuery += fmt.Sprintf(` AND "category" = '%s'`, params.MerchantCategory)
	}

	getMerchantQueryJustWithFilter := getMerchantQuery

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
		if err := rows.Scan(&merchant.ID, &merchant.Name, &merchant.Category, &merchant.ImageURL, &merchant.Location.Lat, &merchant.Location.Long, &merchant.CreatedAt); err != nil {
			return nil, metaData, err
		}
		listMerchant = append(listMerchant, merchant)
	}
	if err := rows.Err(); err != nil {
		return nil, metaData, err
	}

	countQuery := strings.Replace(getMerchantQueryJustWithFilter, "SELECT * FROM", "SELECT count(id) FROM", 1)
	err = r.db.QueryRowxContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, metaData, err
	}

	metaData.Limit = params.Limit
	metaData.Offset = params.Offset
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

func (r *merchantRepository) GetMerchantItem(ctx context.Context, params model.GetMerchantItemParams) (listMerchantItem []model.MerchantItem, meta model.MetaData, err error) {
	listMerchantItem = []model.MerchantItem{}
	var getMerchantItemQuery = `SELECT * FROM "merchantItem" WHERE true`
	var total int = 0
	var metaData = model.MetaData{
		Offset: params.Offset,
		Limit:  params.Limit,
		Total:  0,
	}

	if params.MerchantId != "" {
		getMerchantItemQuery += fmt.Sprintf(` AND "merchantId" = '%s'`, params.MerchantId)
	}

	if params.Name != "" {
		name := "%" + params.Name + "%"
		getMerchantItemQuery += fmt.Sprintf(` AND "name" ILIKE '%s'`, name)
	}

	if params.ItemId != "" {
		getMerchantItemQuery += fmt.Sprintf(` AND "id" = '%s'`, params.ItemId)
	}

	if params.ProductCategory != "" {
		getMerchantItemQuery += fmt.Sprintf(` AND "category" = '%s'`, params.ProductCategory)
	}

	queryWithFilter := getMerchantItemQuery

	if params.CreatedAt != "" {
		if params.CreatedAt != "desc" && params.CreatedAt != "asc" {
			params.CreatedAt = "desc"
		}
		getMerchantItemQuery += fmt.Sprintf(` ORDER BY "createdAt" %s`, params.CreatedAt)
	} else {
		getMerchantItemQuery += ` ORDER BY "createdAt" DESC`
	}

	if params.Limit == 0 {
		params.Limit = 5 // default limit
	}

	getMerchantItemQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, getMerchantItemQuery)
	if err != nil {
		return listMerchantItem, metaData, err
	}

	defer rows.Close()

	// Iterate over the rows and scan each row into a struct
	for rows.Next() {
		var merchantItem model.MerchantItem
		if err := rows.Scan(&merchantItem.ID, &merchantItem.MerchantId, &merchantItem.Name, &merchantItem.Category, &merchantItem.ImageURL, &merchantItem.Price, &merchantItem.CreatedAt); err != nil {
			return listMerchantItem, metaData, err
		}
		listMerchantItem = append(listMerchantItem, merchantItem)
	}
	if err := rows.Err(); err != nil {
		return listMerchantItem, metaData, err
	}

	countQuery := strings.Replace(queryWithFilter, "SELECT * FROM", "SELECT count(id) FROM", 1)
	err = r.db.QueryRowxContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return listMerchantItem, metaData, err
	}

	metaData.Limit = params.Limit
	metaData.Offset = params.Offset
	metaData.Total = total

	return listMerchantItem, metaData, nil
}
