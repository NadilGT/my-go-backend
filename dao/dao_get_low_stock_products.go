package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB_GetTopLowStockProducts returns the top N products with the lowest stock quantity (excluding deleted)
func DB_GetTopLowStockProducts(limit int) ([]dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{"deleted": false, "stockQty": bson.M{"$lt": 10}}
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "stockQty", Value: 1}}) // Ascending order

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []dto.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}
