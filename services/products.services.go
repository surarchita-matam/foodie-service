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


