package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_FindAllBrands() ([]dto.Brand, error) {
	collection := dbConfigs.DATABASE.Collection("Brands")

	ctx := context.Background()
	cursor, err := collection.Find(ctx, map[string]interface{}{})
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
