package service

import (
	"beli-mang/model"
	"beli-mang/repo"

	"github.com/google/uuid"
)

type MerchantService interface {
	CreateMerchant(request model.CreateMerchantRequest) (merchantId string, err error)
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
