package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func DB_FindAllGRNsPaginated(page int, limit int) ([]dto.GRN, int64, error) {
	collection := dbConfigs.DATABASE.Collection("GRNs")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	// Get total count
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip value
	skip := (page - 1) * limit

	// Set up find options with pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var grns []dto.GRN
	if err := cursor.All(ctx, &grns); err != nil {
		return nil, 0, err
	}

	return grns, total, nil
}
