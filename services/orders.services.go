package services

import (
	"fmt"
	"foodie-service/models"
	"foodie-service/types"

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

	return &OrdersService{models: models}
}

func (os *OrdersService) PlaceOrder(order *types.BulkOrdersRequest, userID string) (*types.PurchaseDetails, error) {
	orderID := uuid.New().String()
	// TODO: check for valid coupon code

	totalPrice := 0
	discount := 0
	finalPrice := 0
	products := []types.Product{}
	for _, item := range order.Items {
		product, err := os.models.Products.GetProductByProductId(item.ProductID)
		if err != nil {
			return nil, err
		}
		fmt.Println(product)
		products = append(products, *product)
		totalPrice += int(product.Price) * item.Quantity
		// discount += product.Discount * item.Quantity
		// finalPrice += (product.Price - product.Discount) * item.Quantity
	}
	fmt.Println(products)
	// TODO: apply discount if coupon code is valid
	//    finalPrice = totalPrice % discount /100
	finalPrice = totalPrice - discount
	purchaseDetails := &types.PurchaseDetails{
		OrderID:    orderID,
		Items:      order.Items,
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
