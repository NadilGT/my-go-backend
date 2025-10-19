package dbConfigs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupSalesTTL creates a TTL index on the Sales collection
// Sales will be automatically deleted 24 hours after the created_at timestamp
func SetupSalesTTL() error {
	collection := DATABASE.Collection("Sales")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create TTL index on created_at field
	// expireAfterSeconds: 86400 = 24 hours (24 * 60 * 60 seconds)
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "created_at", Value: 1}, // 1 for ascending order
		},
		Options: options.Index().
			SetExpireAfterSeconds(86400). // 24 hours in seconds
			SetName("sales_ttl_index"),
	}

	indexName, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Error creating TTL index: %v", err)
		return err
	}

	log.Printf("Successfully created TTL index: %s on Sales collection", indexName)
	log.Println("Sales records will be automatically deleted 24 hours after creation")
	return nil
}
