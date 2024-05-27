package service

import (
	"beli-mang/model"
	"beli-mang/repo"
	"context"

	"github.com/google/uuid"
)

type MerchantService interface {
	CreateMerchant(request model.CreateMerchantRequest) (merchantId string, err error)
	GetMerchant(ctx context.Context, params model.GetMerchantParams) (listMerchant []model.Merchant, meta model.MetaData, err error)
	CreateMerchantItem(ctx context.Context, request model.CreateMerchantItemRequest, merchantId string) (itemId string, err error)
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
		ID:        id,
		Name:      request.Name,
		Category:  model.MerchantCategory(request.Category),
		ImageURL:  request.ImageURL,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
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


func (s *merchantSvc) CreateMerchantItem(ctx context.Context, request model.CreateMerchantItemRequest, merchantId string) (itemId string, err error) {

	_, err = s.repo.GetMerchantById(ctx, merchantId)
	if err != nil {
		return "", err
	}

	id := uuid.New()

	merchantItem := model.MerchantItem{
		ID:        id,
		MerchantId: merchantId,
		Name:      request.Name,
		Category:  request.ProductCategory,
		ImageURL:  request.ImageURL,
	}

	err = s.repo.CreateMerchantItem(merchantItem)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}