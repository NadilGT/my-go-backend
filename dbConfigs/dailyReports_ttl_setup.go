package dbConfigs

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupDailyReportsTTL creates a TTL index on the DailyReports collection
// Reports will be automatically deleted by MongoDB when expiresAt date is reached
func SetupDailyReportsTTL() error {
	collection := DATABASE.Collection("DailyReports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create TTL index on expiresAt field
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expiresAt", Value: 1},
		},
		Options: options.Index().
			SetExpireAfterSeconds(0). // Delete immediately when expiresAt is reached
			SetName("dailyReports_ttl_index"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	// Create index on reportDate for efficient queries
	dateIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "reportDate", Value: 1},
		},
		Options: options.Index().
			SetName("dailyReports_date_index").
			SetUnique(false),
	}

	_, err = collection.Indexes().CreateOne(ctx, dateIndexModel)
	if err != nil {
		return err
	}

	// Create compound index on year and month for efficient monthly queries
	monthIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "year", Value: 1},
			{Key: "month", Value: 1},
		},
		Options: options.Index().
			SetName("dailyReports_month_index").
			SetUnique(false),
	}

	_, err = collection.Indexes().CreateOne(ctx, monthIndexModel)
	return err
}
