package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DB_FindAllBrands() ([]dto.Brand, error) {
	collection := dbConfigs.DATABASE.Collection("Brands")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var brands []dto.Brand
	if err := cursor.All(ctx, &brands); err != nil {
		return nil, err
	}

	return brands, nil
}

func DB_FindAllBrandsPaginated(page int, limit int) ([]dto.Brand, int64, error) {
	collection := dbConfigs.DATABASE.Collection("Brands")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	// Get total count
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip value
	skip := (page - 1) * limit

	// Set up find options with pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var brands []dto.Brand
	if err := cursor.All(ctx, &brands); err != nil {
		return nil, 0, err
	}

	return brands, total, nil
}
