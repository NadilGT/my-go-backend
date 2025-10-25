package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_CleanupOrphanedStocks removes stock entries with null or empty batchId
// This fixes data integrity issues where stocks exist without proper batch references
func DB_CleanupOrphanedStocks() (int64, error) {
	stocksCollection := dbConfigs.DATABASE.Collection("Stocks")
	ctx := context.Background()

	// Find and delete stock entries where batchId is null, empty, or doesn't exist
	filter := bson.M{
		"$or": []bson.M{
			{"batchId": nil},
			{"batchId": ""},
			{"batchId": bson.M{"$exists": false}},
		},
	}

	result, err := stocksCollection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup orphaned stocks: %v", err)
	}

	return result.DeletedCount, nil
}

// DB_ValidateStockIntegrity checks for inconsistencies between Products and Stocks collections
// Returns a report of issues found
func DB_ValidateStockIntegrity() (map[string]interface{}, error) {
	productsCollection := dbConfigs.DATABASE.Collection("Products")
	stocksCollection := dbConfigs.DATABASE.Collection("Stocks")
	ctx := context.Background()

	report := make(map[string]interface{})
	issues := []string{}

	// Check 1: Find stocks with null batchId
	nullBatchCount, err := stocksCollection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"batchId": nil},
			{"batchId": ""},
			{"batchId": bson.M{"$exists": false}},
		},
	})
	if err == nil && nullBatchCount > 0 {
		issues = append(issues, fmt.Sprintf("Found %d stock entries with null/empty batchId", nullBatchCount))
	}

	// Check 2: Find stocks referencing non-existent products
	cursor, err := stocksCollection.Distinct(ctx, "productId", bson.M{})
	if err == nil {
		stockProductIds := cursor
		for _, productId := range stockProductIds {
			count, err := productsCollection.CountDocuments(ctx, bson.M{
				"productId": productId,
				"deleted":   false,
			})
			if err == nil && count == 0 {
				issues = append(issues, fmt.Sprintf("Stock references non-existent product: %v", productId))
			}
		}
	}

	// Check 3: Find stocks with batches that don't exist in the product
	// This is more complex and would require iterating through all products
	// For now, we'll skip this check for performance reasons

	report["orphaned_stocks_count"] = nullBatchCount
	report["issues"] = issues
	report["total_issues"] = len(issues)
	report["status"] = "ok"
	if len(issues) > 0 {
		report["status"] = "issues_found"
	}

	return report, nil
}
