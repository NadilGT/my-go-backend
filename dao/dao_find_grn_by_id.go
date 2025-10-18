package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_FindGRNById(grnId string) (*dto.GRN, error) {
	collection := dbConfigs.DATABASE.Collection("GRNs")
	ctx := context.Background()

	filter := bson.M{
		"grnId":   grnId,
		"deleted": false,
	}

	var grn dto.GRN
	err := collection.FindOne(ctx, filter).Decode(&grn)
	if err != nil {
		return nil, err
	}

	return &grn, nil
}
