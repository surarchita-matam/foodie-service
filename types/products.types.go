package types

type Image struct {
	Thumbnail string `json:"thumbnail"`
	Mobile    string `json:"mobile"`
	Tablet    string `json:"tablet"`
	Desktop   string `json:"desktop"`
}

type Product struct {
	ProductID  string    `json:"productId" validate:"required"`
	Image      Image     `json:"image" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Category   string    `json:"category" validate:"required"`
	Price      float64   `json:"price" validate:"required"`
}

type BulkProductsRequest struct {
	Products []Product `json:"products" validate:"required"`
}