package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB_SearchCategories searches for categories by name (case-insensitive) with optional limit
func DB_SearchCategories(searchTerm string, limit int) ([]dto.Category, error) {
	collection := dbConfigs.DATABASE.Collection("Categories")
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

	var categories []dto.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}
