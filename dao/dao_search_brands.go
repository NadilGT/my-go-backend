package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB_SearchBrands searches for brands by name (case-insensitive) with optional limit
func DB_SearchBrands(searchTerm string, limit int) ([]dto.Brand, error) {
	collection := dbConfigs.DATABASE.Collection("Brands")
	ctx := context.Background()

	filter := bson.M{
		"deleted": false,
		"name":    bson.M{"$regex": searchTerm, "$options": "i"},
	}

	// Set up find options with limit
	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	findOptions.SetSort(bson.M{"name": 1}) // Sort by name alphabetically

	cursor, err := collection.Find(ctx, filter, findOptions)
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
