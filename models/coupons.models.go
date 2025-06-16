package models

import (
	"context"
	"fmt"
	"foodie-service/database"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	couponCollection = "coupons"
)

// Coupon represents a coupon document in the database
type Coupon struct {
	ID          string   `json:"_id" bson:"_id"`
	Code        string   `json:"code" bson:"code" unique:"true"`
	FileList    []string `json:"fileList" bson:"fileList"`
	Appearances int      `json:"appearances" bson:"appearances"`
}

type TempCoupon struct {
	Code        string   `json:"code" bson:"code"`
	Appearances []string `json:"appearances" bson:"appearances"`
}

type CouponModel struct {
	dbp *database.Mongo
	dbs *database.Mongo
}

func NewCouponModel(mongoClientPrimary *database.Mongo, mongoClientSecondary *database.Mongo) *CouponModel {
	return &CouponModel{dbp: mongoClientPrimary, dbs: mongoClientSecondary}
}

func (m *CouponModel) InitCollection(ctx context.Context) error {
	db := m.dbp.MongoClient.Database("foodie")
	_, err := db.Collection(couponCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"code": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}
	return nil
}

func (m *CouponModel) CollectionExists(ctx context.Context) (bool, error) {
	db := m.dbp.MongoClient.Database("foodie")
	count, err := db.Collection(couponCollection).CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, fmt.Errorf("failed to check collection: %v", err)
	}
	return count > 0, nil
}

func (m *CouponModel) BulkUpsertCoupons(ctx context.Context, codes []string) error {
	if len(codes) == 0 {
		return nil
	}

	startTime := time.Now()
	totalRecords := len(codes)
	log.Printf("Starting bulk upsert of %d records...", totalRecords)

	const batchSize = 1000
	db := m.dbp.MongoClient.Database("foodie")
	collection := db.Collection(couponCollection)

	for i := 0; i < len(codes); i += batchSize {
		batchStartTime := time.Now()
		end := i + batchSize
		if end > len(codes) {
			end = len(codes)
		}

		models := make([]mongo.WriteModel, end-i)
		for j, code := range codes[i:end] {
			update := mongo.NewUpdateOneModel()
			update.SetFilter(bson.M{"code": code})
			update.SetUpdate(bson.M{"$inc": bson.M{"appearances": 1}})
			update.SetUpsert(true)
			models[j] = update
		}

		opts := options.BulkWrite().SetOrdered(false)
		_, err := collection.BulkWrite(ctx, models, opts)
		if err != nil {
			return fmt.Errorf("bulk upsert failed at batch %d-%d: %v", i, end, err)
		}

		progress := float64(end) / float64(totalRecords) * 100
		batchDuration := time.Since(batchStartTime)
		log.Printf("Processed batch %d-%d (%.2f%%) in %v", i, end, progress, batchDuration)
	}

	totalDuration := time.Since(startTime)
	log.Printf("Completed bulk upsert of %d records in %v", totalRecords, totalDuration)
	return nil
}

func (m *CouponModel) OptimizedBulkInsert(ctx context.Context, coupons []Coupon) error {
	if len(coupons) == 0 {
		return nil
	}

	startTime := time.Now()
	totalRecords := len(coupons)
	log.Printf("Starting optimized bulk insert of %d records...", totalRecords)

	const batchSize = 10000 
	db := m.dbp.MongoClient.Database("foodie")
	collection := db.Collection(couponCollection)

	documents := make([]interface{}, 0, batchSize)
	insertedCount := 0
	totalBatches := (totalRecords + batchSize - 1) / batchSize
	currentBatchNum := 0
	lastProgressUpdate := time.Now()
	lastProgressPercentage := 0.0

	for _, coupon := range coupons { 
	    coupon.ID = primitive.NewObjectID().Hex()
        documents = append(documents, coupon) 

        insertedCount++

        if len(documents) >= batchSize || insertedCount == totalRecords {
            batchStartTime := time.Now()
            currentBatchNum++

            if len(documents) > 0 { 
                opts := options.InsertMany().SetOrdered(false) 
                var insertErr error
                for retries := 0; retries < 3; retries++ { 
                    _, insertErr = collection.InsertMany(ctx, documents, opts)
                    if insertErr == nil {
                        break 
                    }
                    log.Printf("Retry %d for batch %d/%d (%d docs) due to error: %v",
                        retries+1, currentBatchNum, totalBatches, len(documents), insertErr)
                    time.Sleep(time.Second * time.Duration(retries+1))
                }
                if insertErr != nil {
                    return fmt.Errorf("bulk insert failed at batch %d/%d after retries (%d documents): %v",
                        currentBatchNum, totalBatches, len(documents), insertErr)
                }
            }

            currentProgressPercentage := float64(insertedCount) / float64(totalRecords) * 100
            if currentProgressPercentage-lastProgressPercentage >= 5.0 || time.Since(lastProgressUpdate) >= 5*time.Second {
                batchDuration := time.Since(batchStartTime)
                log.Printf("Inserted batch %d/%d (%.2f%%) - %d documents in %v",
                    currentBatchNum, totalBatches, currentProgressPercentage, len(documents), batchDuration)
                lastProgressPercentage = currentProgressPercentage
                lastProgressUpdate = time.Now()
            }

            documents = documents[:0]
        }
    }

	totalDuration := time.Since(startTime)
	log.Printf("Completed optimized bulk insert of %d codes in %v", insertedCount, totalDuration)
	return nil
}

func (m *CouponModel) GetCouponCount(ctx context.Context, code string) (int, error) {
	var coupon Coupon
	db := m.dbp.MongoClient.Database("foodie")
	err := db.Collection(couponCollection).FindOne(ctx, bson.M{"code": code}).Decode(&coupon)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get coupon: %v", err)
	}
	return coupon.Appearances, nil
}

func (m *CouponModel) ValidateCoupon(ctx context.Context, code string) (bool, error) {
	var coupon Coupon
	db := m.dbp.MongoClient.Database("foodie")
	err := db.Collection(couponCollection).FindOne(ctx, bson.M{"code": code}).Decode(&coupon)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, fmt.Errorf("failed to get coupon: %v", err)
	}
	return coupon.Appearances >= 2, nil
}

func (m *CouponModel) FetchCoupons() ([]Coupon, error) {
	var coupons []Coupon
	db := m.dbp.MongoClient.Database("foodie")
	cursor, err := db.Collection(couponCollection).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch coupons: %v", err)
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &coupons)
	if err != nil {
		return nil, fmt.Errorf("failed to decode coupons: %v", err)
	}
	return coupons, nil
}
