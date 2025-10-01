package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_FindProductsBySupplierID(supplierId string) ([]dto.SupplierProduct, error) {
	collection := dbConfigs.DATABASE.Collection("SupplierProducts")
	ctx := context.Background()

	cursor, err := collection.Find(ctx, bson.M{"supplierId": supplierId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assignments []dto.SupplierProduct
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, err
	}

	return assignments, nil
}
