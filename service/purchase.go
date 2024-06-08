package service

import (
	"beli-mang/model"
	cerr "beli-mang/pkg/customErr"
	"beli-mang/pkg/panics"
	"beli-mang/repo"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math"
	"net/http"
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
	totalPrice := 0
	var merchantIds []uuid.UUID
	var itemIds []uuid.UUID
	var detail model.OrderDetail
	var merchantIDStartingPoint uuid.UUID
	mapItemQuantity := make(map[uuid.UUID]int)
	for _, order := range request.Orders {
		merchantId, err := uuid.Parse(order.MerchantId)
		if err != nil {
			return response, cerr.New(http.StatusNotFound, "bad merchant id")
		}
		if order.IsStartingPoint {
			if merchantIDStartingPoint != uuid.Nil {
				return response, cerr.New(http.StatusBadRequest, "multiple starting points")
			}
			merchantIDStartingPoint = merchantId
		}

		merchantIds = append(merchantIds, merchantId)
		for _, item := range order.Items {
			itemID, err := uuid.Parse(item.ItemId)
			if err != nil {
				return response, cerr.New(http.StatusNotFound, "bad item id")
			}
			mapItemQuantity[itemID] = item.Quantity
			itemIds = append(itemIds, itemID)
		}
	}

	if merchantIDStartingPoint == uuid.Nil {
		return response, cerr.New(http.StatusBadRequest, "not have starting point")
	}

	// use go routine to get data concurrently
	var mapMerchant map[uuid.UUID]model.Merchant
	var mapItems map[uuid.UUID]model.Item
	var wg sync.WaitGroup

	wg.Add(1)
	go panics.CaptureGoroutine(func() {
		var errWg error
		defer wg.Done()
		// get merchantsLocation
		mapMerchant, errWg = s.merchantRepo.GetMerchantMapByIds(ctx, merchantIds)
		if errWg != nil {
			s.logger.Error(logPrefix+"failed to get merchant map", zap.Error(errWg))
		}
	}, func() {})

	wg.Add(1)
	go panics.CaptureGoroutine(func() {
		var errWg error
		defer wg.Done()
		// get itemIds Batch
		mapItems, errWg = s.merchantRepo.GetMerchantItemMapByIds(ctx, itemIds)
		if errWg != nil {
			s.logger.Error(logPrefix+"failed to get merchant item map", zap.Error(errWg))
		}
	}, func() {})

	wg.Wait()

	if len(mapItems) == 0 || len(mapMerchant) == 0 {
		return response, cerr.New(http.StatusBadRequest, "invalid items/merchants request")
	}

	for _, item := range mapItems {
		totalPrice += item.Price * mapItemQuantity[item.Id]
	}

	var merchantNames []string
	var merchantCategories []model.MerchantCategory
	var ItemsName []string
	for _, order := range request.Orders {
		var boughtItems []model.BoughtItem
		for _, item := range order.Items {
			itemID, _ := uuid.Parse(item.ItemId)
			bItem := mapItems[itemID].ItemToBoughtItem()
			ItemsName = append(ItemsName, bItem.Name)
			bItem.Quantity = item.Quantity
			boughtItems = append(boughtItems, bItem)
		}
		merchantId, _ := uuid.Parse(order.MerchantId)
		merchant := mapMerchant[merchantId]
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

	// calculate distance by tsp
	end := model.Point{
		Lat: request.UserLocation.Lat,
		Lon: request.UserLocation.Long,
	}
	// compose point from merchantLocation
	points := make([]model.Point, len(mapMerchant))
	for _, merchant := range mapMerchant {
		if merchant.ID == merchantIDStartingPoint {
			points = append([]model.Point{{
				Lat: merchant.Location.Lat,
				Lon: merchant.Location.Long,
			}}, points...)
			continue
		}
		points = append(points, model.Point{
			Lat: merchant.Location.Lat,
			Lon: merchant.Location.Long,
		})

	}

	// check first distance, points will always not null.
	p1 := points[0]
	p2 := end
	if distance := haversineDistance(p1.Lat, p1.Lon, p2.Lat, p2.Lon); distance > model.MaxDistanceFromStartingPoint {
		return response, cerr.New(http.StatusBadRequest, fmt.Sprintf("far from merchant starting point with distance %.2f km", distance))
	}
	estTime := EstimateDeliveryTimeTSP(points, end)

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
