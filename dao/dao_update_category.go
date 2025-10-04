package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_UpdateCategory(ctx context.Context, categoryId string, name string, updatedAt time.Time) error {
	collection := dbConfigs.DATABASE.Collection("Categories")

	filter := bson.M{"categoryId": categoryId}

	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"updated_at": updatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("not_found")
	}

	return nil
}
