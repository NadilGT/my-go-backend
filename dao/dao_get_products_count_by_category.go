package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_GetProductsCountByCategory returns the count of products for a specific category (non-deleted only)
func DB_GetProductsCountByCategory(categoryId string) (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{
		"deleted":    false,
		"categoryId": categoryId,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
