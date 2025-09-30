package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_FindAllSuppliers() ([]dto.Supplier, error) {
	collection := dbConfigs.DATABASE.Collection("Suppliers")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	cursor, err := collection.Find(ctx, filter)
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
