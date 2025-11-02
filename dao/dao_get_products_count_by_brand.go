package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_GetProductsCountByBrand returns the count of products for a specific brand (non-deleted only)
func DB_GetProductsCountByBrand(brandId string) (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{
		"deleted": false,
		"brandId": brandId,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
