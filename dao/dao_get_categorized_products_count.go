package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_GetCategorizedProductsCount returns the count of products that have a categoryId assigned and are not deleted
func DB_GetCategorizedProductsCount() (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{
		"deleted":    false,
		"categoryId": bson.M{"$exists": true, "$ne": ""},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
