package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_DeleteSupplierByID(id string) error {
	collection := dbConfigs.DATABASE.Collection("Suppliers")

	result, err := collection.DeleteOne(context.Background(), bson.M{"supplierId": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return nil
	}

	return nil
}
