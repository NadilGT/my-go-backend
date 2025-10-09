package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DB_FindAllProducts() ([]dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	cursor, err := collection.Find(ctx, filter)
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

func DB_FindAllProductsPaginated(page, limit int) ([]dto.Product, int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	// Calculate skip value
	skip := (page - 1) * limit

	// Count total documents matching the filter
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Set up find options with pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	// Optional: Add sorting for consistent results
	findOptions.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []dto.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
