package service

import (
	"beli-mang/model"
	cerr "beli-mang/pkg/customErr"
	"beli-mang/repo"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type MerchantService interface {
	CreateMerchant(request model.CreateMerchantRequest) (merchantId string, err error)
	GetMerchant(ctx context.Context, params model.GetMerchantParams) (listMerchant []model.Merchant, meta model.MetaData, err error)
	CreateMerchantItem(ctx context.Context, request model.CreateMerchantItemRequest, merchantId uuid.UUID) (itemId string, err error)
	GetMerchantItem(ctx context.Context, merchantId uuid.UUID, params model.GetMerchantItemParams) (listMerchant []model.MerchantItem, meta model.MetaData, err error)
}

type merchantSvc struct {
	repo repo.MerchantRepository
}

func NewMerchantService(r repo.MerchantRepository) MerchantService {
	return &merchantSvc{
		repo: r,
	}
}

func (s *merchantSvc) CreateMerchant(request model.CreateMerchantRequest) (merchantId string, err error) {
	id := uuid.New()

	merchant := model.Merchant{
		ID:       id,
		Name:     request.Name,
		Category: model.MerchantCategory(request.Category),
		ImageURL: request.ImageURL,
		Location: request.Location,
	}

	err = s.repo.CreateMerchant(merchant)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *merchantSvc) GetMerchant(ctx context.Context, params model.GetMerchantParams) (listMerchant []model.Merchant, meta model.MetaData, err error) {
	listMerchant, meta, err = s.repo.GetMerchant(ctx, params)
	if err != nil {
		return
	}

	return listMerchant, meta, nil
}

func (s *merchantSvc) CreateMerchantItem(ctx context.Context, request model.CreateMerchantItemRequest, merchantId uuid.UUID) (itemId string, err error) {

	_, err = s.repo.GetMerchantById(ctx, merchantId)
	if err != nil {
		return "", cerr.New(http.StatusNotFound, err.Error())
	}

	id := uuid.New()

	merchantItem := model.MerchantItem{
		ID:         id,
		MerchantId: merchantId,
		Name:       request.Name,
		Category:   request.ProductCategory,
		ImageURL:   request.ImageURL,
		Price:      request.Price,
		CreatedAt:  time.Now(),
	}

	err = s.repo.CreateMerchantItem(merchantItem)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *merchantSvc) GetMerchantItem(ctx context.Context, merchantId uuid.UUID, params model.GetMerchantItemParams) (listMerchantItem []model.MerchantItem, meta model.MetaData, err error) {
	listMerchantItem = []model.MerchantItem{}
	_, err = s.repo.GetMerchantById(ctx, merchantId)
	if err != nil {
		return listMerchantItem, meta, cerr.New(http.StatusNotFound, err.Error())
	}

	params.MerchantId = merchantId.String()
	listMerchantItem, meta, err = s.repo.GetMerchantItem(ctx, params)
	if err != nil {
		return listMerchantItem, meta, cerr.New(http.StatusInternalServerError, err.Error())
	}

	return listMerchantItem, meta, nil
}
