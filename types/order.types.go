package types

import "time"

type Order struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required"`
}

type BulkOrdersRequest struct {
	Items      []Order `json:"items" validate:"required"`
	CouponCode string  `json:"couponCode"`
}

type PurchaseDetails struct {
	OrderID    string    `json:"orderId"`
	Items      []Order   `json:"items"`
	Products   []Product `json:"products"`
	TotalPrice float64      `json:"totalPrice"`
	Discount   float64       `json:"discount"`
	FinalPrice float64       `json:"finalPrice"`
	CouponCode string    `json:"couponCode"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
