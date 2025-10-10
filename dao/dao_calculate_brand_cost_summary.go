package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_GetBrandCostSummary(brandId string) (float64, float64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	matchStage := bson.D{{Key: "$match", Value: bson.M{
		"deleted": false,
		"brandId": brandId,
	}}}

	addFieldsStage := bson.D{{Key: "$addFields", Value: bson.M{
		"totalCost":    bson.M{"$multiply": []string{"$costPrice", "$stockQty"}},
		"expectedCost": bson.M{"$multiply": []string{"$sellingPrice", "$stockQty"}},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.M{
		"_id":           "$brandId",
		"total_cost":    bson.M{"$sum": "$totalCost"},
		"expected_cost": bson.M{"$sum": "$expectedCost"},
	}}}

	cursor, err := collection.Aggregate(ctx, bson.A{matchStage, addFieldsStage, groupStage})
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0, 0, err
	}

	if len(results) > 0 {
		totalCost, _ := results[0]["total_cost"].(float64)
		expectedCost, _ := results[0]["expected_cost"].(float64)
		return totalCost, expectedCost, nil
	}

	return 0, 0, nil
}
