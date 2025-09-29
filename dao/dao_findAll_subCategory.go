package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_FindAllSubCategory() ([]dto.SubCategory, error) {
	collection := dbConfigs.DATABASE.Collection("SubCategories")

	ctx := context.Background()
	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subCategory []dto.SubCategory
	if err := cursor.All(ctx, &subCategory); err != nil {
		return nil, err
	}

	return subCategory, nil
}
