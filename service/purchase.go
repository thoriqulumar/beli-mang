package service

import (
	"beli-mang/model"
	"beli-mang/pkg/panics"
	"beli-mang/repo"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math"
	"strings"
	"sync"
)

type PurchaseService interface {
	EstimateOrders(ctx context.Context, request model.EstimateOrdersRequest) (response model.EstimateOrdersResponse, err error)
	ConfirmOrder(ctx context.Context, request model.ConfirmOrderRequest) (response model.ConfirmOrderResponse, err error)
	GetUserOrders(ctx context.Context, request model.UserOrdersParams) (response model.GetUserOrdersResponse, err error)
	GetNearbyMerchant(ctx context.Context, params model.GetMerchantParams, lat, long string) (listMerchant []model.GetNearbyMerchantData, meta model.MetaData, err error)
}

type purchaseSvc struct {
	orderRepo    repo.OrderRepository
	merchantRepo repo.MerchantRepository
	logger       *zap.Logger
}

func NewPurchaseService(orderRepo repo.OrderRepository, merchantRepo repo.MerchantRepository, logger *zap.Logger) PurchaseService {
	return &purchaseSvc{
		orderRepo:    orderRepo,
		merchantRepo: merchantRepo,
		logger:       logger,
	}
}

/*
try use one table transaction
*/

func (s *purchaseSvc) EstimateOrders(ctx context.Context, request model.EstimateOrdersRequest) (response model.EstimateOrdersResponse, err error) {
	logPrefix := "[purchase] EstimateOrders"
	// calculate distance by tsp
	end := model.Point{
		Lat: request.UserLocation.Lat,
		Lon: request.UserLocation.Long,
	}
	totalPrice := 0
	var merchantIds []uuid.UUID
	var itemIds []uuid.UUID
	var detail model.OrderDetail
	mapItemQuantity := make(map[uuid.UUID]int)
	for _, order := range request.Orders {
		merchantIds = append(merchantIds, order.MerchantId)
		for _, item := range order.Items {
			mapItemQuantity[item.ItemId] = item.Quantity
			itemIds = append(itemIds, item.ItemId)
		}
	}

	// use go routine to get data concurrently
	var mapMerchant map[uuid.UUID]model.Merchant
	var mapItems map[uuid.UUID]model.Item
	var errWg error
	var wg sync.WaitGroup

	go panics.CaptureGoroutine(func() {
		wg.Add(1)
		// get merchantsLocation
		mapMerchant, errWg = s.merchantRepo.GetMerchantMapByIds(ctx, merchantIds)
		if errWg != nil {
			s.logger.Error(logPrefix+"failed to get merchant map", zap.Error(errWg))
		}
	}, func() {})

	go panics.CaptureGoroutine(func() {
		wg.Add(1)
		// get itemIds Batch
		mapItems, errWg = s.merchantRepo.GetMerchantItemMapByIds(ctx, itemIds)
		if errWg != nil {
			s.logger.Error(logPrefix+"failed to get merchant item map", zap.Error(errWg))
		}
	}, func() {})

	wg.Wait()

	for _, item := range mapItems {
		totalPrice += item.Price * mapItemQuantity[item.Id]
	}

	var merchantNames []string
	var merchantCategories []model.MerchantCategory
	var ItemsName []string
	for _, order := range request.Orders {
		var boughtItems []model.BoughtItem
		for _, item := range order.Items {
			bItem := mapItems[item.ItemId].ItemToBoughtItem()
			ItemsName = append(ItemsName, bItem.Name)
			bItem.Quantity = item.Quantity
			boughtItems = append(boughtItems, bItem)
		}
		merchant := mapMerchant[order.MerchantId]
		merchantNames = append(merchantNames, merchant.Name)
		merchantCategories = append(merchantCategories, merchant.Category)
		detail = append(detail, model.OrderData{
			Merchant:        merchant,
			IsStartingPoint: order.IsStartingPoint,
			Items:           boughtItems,
		})
	}
	detailRaw, err := json.Marshal(detail)
	if err != nil {
		return response, err
	}

	// compose point from merchantLocation
	estTime := EstimateDeliveryTimeTSP([]model.Point{}, end)

	// tx start
	tx, err := s.orderRepo.BeginTx(ctx)
	if err != nil {
		return response, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	// submit order draft
	orderId := uuid.New()
	orderData := model.Order{
		OrderID:            orderId,
		OrderStatus:        model.OrderStatusDraft,
		DetailRaw:          detailRaw,
		Detail:             detail,
		MerchantIDs:        merchantIds,
		JoinedMerchantName: strings.Join(merchantNames, ";"),
		MerchantCategories: merchantCategories,
		JoinedItemsName:    strings.Join(ItemsName, ";"),
		UserID:             request.UserId, // buyerId
		UserLatitude:       request.UserLocation.Lat,
		UserLongitude:      request.UserLocation.Long,
	}
	_, err = s.orderRepo.Create(ctx, tx, orderData)
	if err != nil {
		return response, err
	}

	// submit calculation
	calculatedData := model.CalculatedEstimate{
		TotalPrice:                     totalPrice,
		EstimatedDeliveryTimeInMinutes: int(math.Round(estTime.Minutes())),
		CalculatedEstimateId:           uuid.New(),
		OrderId:                        orderId,
	}
	_, err = s.orderRepo.InsertCalculation(ctx, tx, calculatedData)
	if err != nil {
		return response, err
	}

	// submit order
	return model.EstimateOrdersResponse{
		TotalPrice:                     calculatedData.TotalPrice,
		EstimatedDeliveryTimeInMinutes: calculatedData.EstimatedDeliveryTimeInMinutes,
		CalculatedEstimateId:           calculatedData.CalculatedEstimateId,
	}, nil
}

func (s *purchaseSvc) ConfirmOrder(ctx context.Context, request model.ConfirmOrderRequest) (response model.ConfirmOrderResponse, err error) {
	calculatedData, err := s.orderRepo.GetCalculatedEstimateById(ctx, request.CalculatedEstimateId)
	if err != nil {
		return response, err
	}
	err = s.orderRepo.UpdateStatus(ctx, calculatedData.OrderId, model.OrderStatusCreated)
	if err != nil {
		return response, err
	}
	response.OrderId = calculatedData.OrderId
	return response, nil
}

func (s *purchaseSvc) GetUserOrders(ctx context.Context, request model.UserOrdersParams) (response model.GetUserOrdersResponse, err error) {
	// get userOrder
	request.Status = model.OrderStatusCreated
	listData, err := s.orderRepo.GetUserOrders(ctx, request)
	if err != nil {
		return response, err
	}

	for _, order := range listData {
		response = append(response, order.ToUserOrderData())
	}

	return
}


func (s *purchaseSvc) GetNearbyMerchant(ctx context.Context, params model.GetMerchantParams, lat, long string) (listMerchant []model.GetNearbyMerchantData, meta model.MetaData, err error) {
	listMerchant, meta, err = s.orderRepo.GetNearbyMerchant(ctx, params, lat, long)
	if err != nil {
		return
	}

	return listMerchant, meta, nil
}