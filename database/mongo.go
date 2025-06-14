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
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// get mongo uri from env
	opts := options.Client().ApplyURI(config.GetConfig().MONGO_URI).SetServerAPIOptions(serverAPI)

	if readPreference == "primary" {
		opts.SetReadPreference(readpref.Primary())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize common MongoDB client: %w", err)
	}

	// Send a ping to confirm a successful connection
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
