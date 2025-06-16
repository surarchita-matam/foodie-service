package services

import (
	"foodie-service/models"
)

type BaseService struct {
	Products *ProductsService
	Orders   *OrdersService
	Auth     *AuthService
	Coupons  *CouponService
}

var baseService *BaseService

func NewBaseService(models *models.BaseModel) *BaseService {
	if baseService != nil {
		return baseService
	}

	baseService = &BaseService{
		Products: NewProductsService(models),
		Orders:   NewOrdersService(models),
		Auth:     NewAuthService(models),
		Coupons:  NewCouponService(models),
	}
	return baseService
}
