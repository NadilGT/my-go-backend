package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_GetProductCountsForBrands returns a map of brandId -> product count for all brands
func DB_GetProductCountsForBrands(brandIds []string) (map[string]int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	// Create a pipeline to count products grouped by brandId
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"deleted": false,
				"brandId": bson.M{"$in": brandIds},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$brandId",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Parse results into a map
	result := make(map[string]int64)
	for cursor.Next(ctx) {
		var doc struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		result[doc.ID] = doc.Count
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
