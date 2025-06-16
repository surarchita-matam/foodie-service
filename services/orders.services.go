package services

import (
	"context"
	"fmt"
	"foodie-service/models"
	"foodie-service/types"
    "math"
	"github.com/google/uuid"
)

type OrdersService struct {
	models *models.BaseModel
}

var ordersService *OrdersService

func NewOrdersService(models *models.BaseModel) *OrdersService {
	if ordersService != nil {
		return ordersService
	}

	ordersService= &OrdersService{
		models: models,
	}
	return ordersService
}

func (os *OrdersService) PlaceOrder(order *types.BulkOrdersRequest, userID string) (*types.PurchaseDetails, error) {
	orderID := uuid.New().String()
	totalPrice := 0.0
	discount := 0.0
	finalPrice := 0.0
	products := []types.Product{}

	// Calculate total price and get products
	for _, item := range order.Items {
		product, err := os.models.Products.GetProductByProductId(item.ProductID)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
		totalPrice += product.Price * float64(item.Quantity)
	}

	// Validate and apply coupon code if provided
	if order.CouponCode != "" {
		isValid, err := os.models.Coupons.ValidateCoupon(context.Background(), order.CouponCode)
		if err != nil {
			return nil, err
		}
		if isValid {
			discount = math.Round(float64(totalPrice)*0.10 * 100) / 100
		} else {
			return nil, fmt.Errorf("invalid coupon code: %s", order.CouponCode)
		}
	}

	finalPrice = totalPrice - discount

	purchaseDetails := &types.PurchaseDetails{
		OrderID:    orderID,
		Items:      order.Items,
		Products:   products,
		TotalPrice: totalPrice,
		Discount:   discount,
		FinalPrice: finalPrice,
		CouponCode: order.CouponCode,
	}

	purchaseDetails, err := os.models.Orders.InsertOrder(purchaseDetails, userID)
	if err != nil {
		return nil, err
	}
	purchaseDetails.Products = products

	return purchaseDetails, nil
}

func (os *OrdersService) GetPreviousOrders(userID string, limit, offset int)(*[]types.PurchaseDetails, error) {
  orderSchemas, err := os.models.Orders.GetOrders(userID, limit, offset)
  if err != nil {
    return nil, err
  }

  purchaseDetails := []types.PurchaseDetails{}
  for _, orderSchema := range orderSchemas {
	products := []types.Product{}
	for _, item := range orderSchema.Items {
		product, err := os.models.Products.GetProductByProductId(item.ProductID)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}
    purchaseDetails = append(purchaseDetails, types.PurchaseDetails{
      OrderID: orderSchema.OrderID,
      Items: orderSchema.Items,
      TotalPrice: orderSchema.TotalPrice,
      Discount: orderSchema.Discount,
      FinalPrice: orderSchema.FinalPrice,
	  CouponCode: orderSchema.CouponCode,
	  Products: products,
    })
  }

  return &purchaseDetails, nil
}
