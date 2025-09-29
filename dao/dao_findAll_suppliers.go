package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_FindAllSuppliers() ([]dto.Supplier, error) {
	collection := dbConfigs.DATABASE.Collection("Suppliers")

	ctx := context.Background()
	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var suppliers []dto.Supplier
	if err := cursor.All(ctx, &suppliers); err != nil {
		return nil, err
	}

	return suppliers, nil
}
