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
	TotalPrice int       `json:"totalPrice"`
	Discount   int       `json:"discount"`
	FinalPrice int       `json:"finalPrice"`
	CouponCode string    `json:"couponCode"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
