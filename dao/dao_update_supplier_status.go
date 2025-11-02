package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_UpdateSupplierStatus(ctx context.Context, supplierId string, status string) error {
	collection := dbConfigs.DATABASE.Collection("Suppliers")
	filter := bson.M{"supplierId": supplierId}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
