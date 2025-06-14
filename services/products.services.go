package services

import (
	"foodie-service/models"
)

type ProductsService struct {
	models *models.BaseModel
}

func NewProductsService(models *models.BaseModel) *ProductsService {
	return &ProductsService{
		models: models,
	}
}

// Add your product service methods here
// Example:
// func (s *ProductsService) GetProduct(id string) (*models.Product, error) {
//     return s.model.GetProduct(id)
// }

