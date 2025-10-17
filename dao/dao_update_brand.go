package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_UpdateBrand(ctx context.Context, brandId string, name string, categoryId string, updatedAt time.Time) error {
	collection := dbConfigs.DATABASE.Collection("Brands")

	filter := bson.M{"brandId": brandId}

	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"categoryId": categoryId,
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
