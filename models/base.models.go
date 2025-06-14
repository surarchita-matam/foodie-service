package models

import (
	"foodie-service/database"
)

type BaseModel struct {
	Products *ProductsModel
	Orders   *OrdersModel
	Auth     *AuthModel
}

var baseModel *BaseModel

func NewBaseModel(mongoClientPrimary *database.Mongo, mongoClientSecondary *database.Mongo) *BaseModel {
	if baseModel != nil {
		return baseModel
	}

	baseModel = &BaseModel{
		Products: NewProductsModel(mongoClientPrimary, mongoClientSecondary),
		Orders:   NewOrdersModel(mongoClientPrimary, mongoClientSecondary),
		Auth:     NewAuthModel(mongoClientPrimary, mongoClientSecondary),
	}
	return baseModel
}
