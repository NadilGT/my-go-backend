package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DB_CheckGRNExists checks if a GRN exists and is not deleted
func DB_CheckGRNExists(grnId string) (bool, error) {
	collection := dbConfigs.DATABASE.Collection("GRNs")
	ctx := context.Background()

	filter := bson.M{
		"grnId":   grnId,
		"deleted": false,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// DB_UpdateGRNStatus updates the status of a GRN
func DB_UpdateGRNStatus(grnId string, status string, updatedAt time.Time) error {
	collection := dbConfigs.DATABASE.Collection("GRNs")
	ctx := context.Background()

	filter := bson.M{
		"grnId":   grnId,
		"deleted": false,
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": updatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
