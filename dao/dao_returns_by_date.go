package dao

import (
	"context"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ...existing code...

func GetReturnsByDateRange(ctx context.Context, start, end time.Time) ([]dto.ReturnDTO, error) {
	filter := bson.M{
		"createdat": bson.M{
			"$gte": start.Format(time.RFC3339),
			"$lte": end.Format(time.RFC3339),
		},
	}
	cursor, err := ReturnsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []dto.ReturnDTO
	for cursor.Next(ctx) {
		var r dto.ReturnDTO
		if err := cursor.Decode(&r); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}
