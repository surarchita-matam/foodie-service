package controllers

import (
	"foodie-service/models"
	"foodie-service/services"
)

type BaseController struct {
	ProductsController *ProductsController
	OrdersController   *OrdersController
	AuthController     *AuthController
}

var baseController *BaseController

func NewBaseController(services *services.BaseService, models *models.BaseModel) *BaseController {
	if baseController != nil {
		return baseController
	}

	baseController = &BaseController{
		ProductsController: NewProductsController(services, models),
		OrdersController:   NewOrdersController(services, models),
		AuthController:     NewAuthController(services, models),
	}
	return baseController
}

func GetController() *BaseController {
	return baseController
}
