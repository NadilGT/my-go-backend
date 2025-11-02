package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_GetSupplierStatusCounts(ctx context.Context) (active int64, inactive int64, err error) {
	collection := dbConfigs.DATABASE.Collection("Suppliers")
	activeCount, err := collection.CountDocuments(ctx, bson.M{"status": "active", "deleted": false})
	if err != nil {
		return 0, 0, err
	}
	inactiveCount, err := collection.CountDocuments(ctx, bson.M{"status": "inactive", "deleted": false})
	if err != nil {
		return 0, 0, err
	}
	return activeCount, inactiveCount, nil
}
