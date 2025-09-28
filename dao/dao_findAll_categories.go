package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_FindAllCategories() ([]dto.Category, error) {
	collection := dbConfigs.DATABASE.Collection("Categories")

	ctx := context.Background()
	cursor, err := collection.Find(ctx, map[string]interface{}{})
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
