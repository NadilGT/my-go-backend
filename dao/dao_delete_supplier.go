package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_DeleteSupplierByID(supplierId string) error {

	collection := dbConfigs.DATABASE.Collection("Suppliers")

	filter := bson.M{
		"supplierId": supplierId,
		"deleted":    false,
	}

	update := bson.M{
		"$set": bson.M{"deleted": true},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("Specified Id not found or already deleted!")
	}

	return nil
}
