package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_CalculateTotalStockQuantity calculates the sum of all stockQty in the Products collection
// Returns the total quantity of all products in stock (includes all products, not just those with batches)
func DB_CalculateTotalStockQuantity() (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	// MongoDB aggregation pipeline to sum all stockQty from non-deleted products
	pipeline := []bson.M{
		{
			"$match": bson.M{"deleted": false}, // Only non-deleted products
		},
		{
			"$group": bson.M{
				"_id":           nil,
				"totalStockQty": bson.M{"$sum": "$stockQty"},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		TotalStockQty int64 `bson:"totalStockQty"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	// If no products found, return 0
	if len(result) == 0 {
		return 0, nil
	}

	return result[0].TotalStockQty, nil
}
