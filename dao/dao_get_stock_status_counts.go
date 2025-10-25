package dao

import (
	"context"
	"employee-crud/dbConfigs"

	"go.mongodb.org/mongo-driver/bson"
)

// StockStatusCounts represents the count of stocks by status
type StockStatusCounts struct {
	LowStock     int64 `json:"low_stock"`
	AverageStock int64 `json:"average_stock"`
	GoodStock    int64 `json:"good_stock"`
	Total        int64 `json:"total"`
}

// DB_GetStockStatusCounts returns the count of stocks grouped by status
// This queries the Products collection directly to include all products
func DB_GetStockStatusCounts() (*StockStatusCounts, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	// MongoDB aggregation pipeline to categorize and count products by stock status
	pipeline := []bson.M{
		{
			"$match": bson.M{"deleted": false}, // Only non-deleted products
		},
		{
			"$project": bson.M{
				"status": bson.M{
					"$switch": bson.M{
						"branches": []bson.M{
							{
								"case": bson.M{"$lt": []interface{}{"$stockQty", 10}},
								"then": "Low Stock",
							},
							{
								"case": bson.M{
									"$and": []bson.M{
										{"$gte": []interface{}{"$stockQty", 10}},
										{"$lt": []interface{}{"$stockQty", 25}},
									},
								},
								"then": "Average Stock",
							},
						},
						"default": "Good Stock",
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Initialize counts
	counts := &StockStatusCounts{
		LowStock:     0,
		AverageStock: 0,
		GoodStock:    0,
		Total:        0,
	}

	// Map results to counts
	for _, result := range results {
		switch result.ID {
		case "Low Stock":
			counts.LowStock = result.Count
		case "Average Stock":
			counts.AverageStock = result.Count
		case "Good Stock":
			counts.GoodStock = result.Count
		}
		counts.Total += result.Count
	}

	return counts, nil
}
