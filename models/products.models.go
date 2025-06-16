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

/*
   {
       "id": "1",
       "image": {
           "thumbnail": "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
           "mobile": "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
           "tablet": "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
           "desktop": "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg"
       },
       "name": "Waffle with Berries",
       "category": "Waffle",
       "price": 6.5
   }
*/

type Image struct {
	Thumbnail string `json:"thumbnail" bson:"thumbnail"`
	Mobile    string `json:"mobile" bson:"mobile"`
	Tablet    string `json:"tablet" bson:"tablet"`
	Desktop   string `json:"desktop" bson:"desktop"`
}

type Product struct {
	ID         string    `json:"_id" bson:"_id"`
	ProductID  string    `json:"productId" bson:"productId" unique:"true"`
	Image      Image     `json:"image" bson:"image"`
	Name       string    `json:"name" bson:"name"`
	Category   string    `json:"category" bson:"category"`
	Price      float64   `json:"price" bson:"price"`
	InsertedAt time.Time `json:"insertedAt" bson:"insertedAt"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
}

type ProductsModel struct {
	dbp *database.Mongo
	dbs *database.Mongo
}

func (pm *ProductsModel) createUniqueIndex() error {
	collection := pm.dbp.MongoClient.Database("foodie").Collection("products")

	// Create a unique index on productId
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"productId": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}

func NewProductsModel(dbp *database.Mongo, dbs *database.Mongo) *ProductsModel {
	pm := &ProductsModel{dbp: dbp, dbs: dbs}

	if err := pm.createUniqueIndex(); err != nil {
		panic(fmt.Sprintf("failed to create unique index: %v", err))
	}
	return pm
}

func (pm *ProductsModel) GetProducts(readFromPrimary bool) ([]types.Product, error) {
	var db *database.Mongo

	if readFromPrimary {
		db = pm.dbp
	} else {
		db = pm.dbs
	}

	collection := db.MongoClient.Database("foodie").Collection("products")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []Product
	if err = cursor.All(context.TODO(), &products); err != nil {
		return nil, err
	}
	result := make([]types.Product, len(products))
	for i, p := range products {
		result[i] = types.Product{
			ProductID: p.ProductID,
			Image:     types.Image(p.Image),
			Name:      p.Name,
			Category:  p.Category,
			Price:     p.Price,
		}
	}

	return result, nil
}

func (pm *ProductsModel) GetProductByProductId(id string) (*types.Product, error) {
	db := pm.dbp
	collection := db.MongoClient.Database("foodie").Collection("products")

	filter := bson.M{"productId": id}

	product := &types.Product{}
	err := collection.FindOne(context.TODO(), filter).Decode(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pm *ProductsModel) InsertBulkProducts(products []types.Product) error {
	db := pm.dbp
	collection := db.MongoClient.Database("foodie").Collection("products")

	// Convert products to mongo.WriteModel operations
	bulkWrite := make([]mongo.WriteModel, len(products))
	for i, product := range products {
		modelProduct := Product{
			ID:         primitive.NewObjectID().Hex(),
			ProductID:  product.ProductID,
			Image:      Image(product.Image),
			Name:       product.Name,
			Category:   product.Category,
			Price:      product.Price,
			InsertedAt: time.Now(),
			UpdatedAt:  time.Now(),
		}
		bulkWrite[i] = mongo.NewInsertOneModel().SetDocument(modelProduct)
	}

	_, err := collection.BulkWrite(context.TODO(), bulkWrite)
	if err != nil {
		// Check if error is due to duplicate key
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("duplicate productId found: %v", err)
		}
		return err
	}

	return nil
}
