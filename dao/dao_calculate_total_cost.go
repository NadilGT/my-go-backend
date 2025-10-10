package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_CalculateTotalCostPrice() (float64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	matchStage := bson.D{{Key: "$match", Value: bson.M{"deleted": false}}}

	addFieldsStage := bson.D{{Key: "$addFields", Value: bson.M{
		"totalCost": bson.M{"$multiply": []string{"$costPrice", "$stockQty"}},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.M{
		"_id":         nil,
		"total_value": bson.M{"$sum": "$totalCost"},
	}}}

	cursor, err := collection.Aggregate(ctx, bson.A{matchStage, addFieldsStage, groupStage})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	if len(results) > 0 {
		if total, ok := results[0]["total_value"].(float64); ok {
			return total, nil
		}
	}

	return 0, nil
}
