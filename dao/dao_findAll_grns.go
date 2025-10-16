package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_FindAllGRNs() ([]dto.GRN, error) {
	collection := dbConfigs.DATABASE.Collection("GRNs")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var grns []dto.GRN
	if err := cursor.All(ctx, &grns); err != nil {
		return nil, err
	}

	return grns, nil
}
