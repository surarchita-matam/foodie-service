package database

import (
	"context"
	"fmt"
	"foodie-service/config"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Mongo struct {
	Ctx         context.Context
	Cancel      context.CancelFunc
	MongoClient *mongo.Client
}

var mutex = &sync.Mutex{}

var mongoClientPrimary *Mongo
var mongoClientSecondary *Mongo

func MongoClient(readPreference string) (*Mongo, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if readPreference == "primary" {
		if mongoClientPrimary != nil {
			return mongoClientPrimary, nil
		}
	} else {
		if mongoClientSecondary != nil {
			return mongoClientSecondary, nil
		}
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.GetConfig().MONGO_URI).SetServerAPIOptions(serverAPI)

	if readPreference == "primary" {
		opts.SetReadPreference(readpref.Primary())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(opts)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize common MongoDB client: %w", err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return &Mongo{
		Ctx:         ctx,
		Cancel:      cancel,
		MongoClient: client,
	}, nil
}
