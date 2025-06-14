package models

import (
	"context"
	"foodie-service/database"
	"foodie-service/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderSchema struct {
	ID         string        `json:"_id" bson:"_id"`
	OrderID    string        `json:"orderId" bson:"orderId" unique:"true"`
	UserID     string        `json:"userId" bson:"userId"`
	Items      []types.Order `json:"items" bson:"items"`
	TotalPrice float64       `json:"totalPrice" bson:"totalPrice"`
	Discount   float64       `json:"discount" bson:"discount"`
	FinalPrice float64       `json:"finalPrice" bson:"finalPrice"`
	CouponCode string        `json:"couponCode" bson:"couponCode"`
	InsertedAt time.Time     `json:"insertedAt" bson:"insertedAt"`
	UpdatedAt  time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type OrdersModel struct {
	dbp *database.Mongo
	dbs *database.Mongo
}

func NewOrdersModel(dbp *database.Mongo, dbs *database.Mongo) *OrdersModel {
	return &OrdersModel{dbp: dbp, dbs: dbs}
}

func (om *OrdersModel) GetOrders(readFromPrimary bool) ([]types.Order, error) {
	var db *database.Mongo

	if readFromPrimary {
		db = om.dbp
	} else {
		db = om.dbs
	}

	collection := db.MongoClient.Database("foodie").Collection("orders")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var orders []types.Order
	if err = cursor.All(context.TODO(), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (om *OrdersModel) InsertOrder(order *types.PurchaseDetails, userID string) (*types.PurchaseDetails, error) {
	db := om.dbp
	collection := db.MongoClient.Database("foodie").Collection("orders")

	orderSchema := &OrderSchema{
		ID:         primitive.NewObjectID().Hex(),
		OrderID:    order.OrderID,
		UserID:     userID,
		Items:      order.Items,
		TotalPrice: float64(order.TotalPrice),
		Discount:   float64(order.Discount),
		FinalPrice: float64(order.FinalPrice),
		CouponCode: order.CouponCode,
		InsertedAt: time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), orderSchema)
	if err != nil {
		return nil, err
	}

	return &types.PurchaseDetails{
		OrderID:   orderSchema.OrderID,
		Items:     orderSchema.Items,
		TotalPrice: int(orderSchema.TotalPrice),
		Discount:   int(orderSchema.Discount),
		FinalPrice: int(orderSchema.FinalPrice),
		CouponCode: orderSchema.CouponCode,
		CreatedAt:  orderSchema.InsertedAt,
		UpdatedAt:  orderSchema.UpdatedAt,
	}, nil
}
