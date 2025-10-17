package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_CountTotalGRNs() (int64, error) {
	collection := dbConfigs.DATABASE.Collection("GRNs")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
