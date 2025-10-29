package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_GetTotalProducts returns the total number of non-deleted products in the database
func DB_GetTotalProducts() (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()
	filter := bson.M{"deleted": false}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
