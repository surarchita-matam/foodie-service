package models

import (
	"context"
	"fmt"
	"foodie-service/database"
	"foodie-service/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (om *OrdersModel) createUniqueIndex() error {
	collection := om.dbp.MongoClient.Database("foodie").Collection("orders")

	// Create unique indexes for email and userId separately
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.M{"orderId": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"userId": 1},
			Options: options.Index(),
		},
	}
	// Create all indexes
	for _, model := range indexModels {
		_, err := collection.Indexes().CreateOne(context.TODO(), model)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}
	return nil
}

func NewOrdersModel(dbp *database.Mongo, dbs *database.Mongo) *OrdersModel {
	om := &OrdersModel{dbp: dbp, dbs: dbs}

	if err := om.createUniqueIndex(); err != nil {
		panic(fmt.Sprintf("failed to create unique index: %v", err))
	}
	return om
}

func (om *OrdersModel) GetOrders(userID string, limit, offset int) ([]OrderSchema, error) {
	collection := om.dbp.MongoClient.Database("foodie").Collection("orders")
	findOptions := options.Find().
		SetSort(bson.M{"_id": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := collection.Find(context.TODO(), bson.M{"userId": userID}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var orderSchemas []OrderSchema
	if err = cursor.All(context.TODO(), &orderSchemas); err != nil {
		return nil, err
	}
	// Log for debugging
	fmt.Printf("Total orders for user %s, Requested limit: %d, offset: %d, Got orders: %d\n",
		userID, limit, offset, len(orderSchemas))
	return orderSchemas, nil
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
		OrderID:    orderSchema.OrderID,
		Items:      orderSchema.Items,
		TotalPrice: float64(orderSchema.TotalPrice),
		Discount:   float64(orderSchema.Discount),
		FinalPrice: float64(orderSchema.FinalPrice),
		CouponCode: orderSchema.CouponCode,
		CreatedAt:  orderSchema.InsertedAt,
		UpdatedAt:  orderSchema.UpdatedAt,
	}, nil
}
