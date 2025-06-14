package models

import (
	"context"
	"fmt"
	"foodie-service/database"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSchema struct {
	ID        string    `json:"_id" bson:"_id"`
	UserID    string    `json:"userId" bson:"userId" unique:"true"`
	Email     string    `json:"email" bson:"email" unique:"true"`
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type AuthModel struct {
	dbp *database.Mongo
	dbs *database.Mongo
}

func NewAuthModel(mongoClientPrimary *database.Mongo, mongoClientSecondary *database.Mongo) *AuthModel {
	return &AuthModel{
		dbp: mongoClientPrimary,
		dbs: mongoClientSecondary,
	}
}

func (am *AuthModel) GetUserByEmail(email string) (*UserSchema, error) {
	db := am.dbp
	collection := db.MongoClient.Database("foodie").Collection("users")

	filter := bson.M{"email": email}

	var user *UserSchema
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(user)
	return user, nil
}

func (am *AuthModel) CreateUser(user *UserSchema) (*UserSchema, error) {
	db := am.dbp
	collection := db.MongoClient.Database("foodie").Collection("users")

	user.ID = primitive.NewObjectID().Hex()
	user.UserID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
