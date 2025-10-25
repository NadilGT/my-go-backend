package main

import (
	"context"
	"employee-crud/dao"
	"employee-crud/dbConfigs"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// FixStockIntegrityIssues is a one-time migration script to fix the null batchId issue
// Run this script to clean up existing data
func FixStockIntegrityIssues() error {
	ctx := context.Background()
	stocksCollection := dbConfigs.DATABASE.Collection("Stocks")

	log.Println("=== Starting Stock Integrity Fix ===")

	// Step 1: Count orphaned stocks
	log.Println("\n[Step 1/4] Counting orphaned stocks...")
	orphanedFilter := bson.M{
		"$or": []bson.M{
			{"batchId": nil},
			{"batchId": ""},
			{"batchId": bson.M{"$exists": false}},
		},
	}

	orphanedCount, err := stocksCollection.CountDocuments(ctx, orphanedFilter)
	if err != nil {
		return fmt.Errorf("failed to count orphaned stocks: %v", err)
	}
	log.Printf("Found %d orphaned stock entries with null/empty batchId", orphanedCount)

	// Step 2: Delete orphaned stocks
	if orphanedCount > 0 {
		log.Println("\n[Step 2/4] Deleting orphaned stocks...")
		result, err := stocksCollection.DeleteMany(ctx, orphanedFilter)
		if err != nil {
			return fmt.Errorf("failed to delete orphaned stocks: %v", err)
		}
		log.Printf("Successfully deleted %d orphaned stock entries", result.DeletedCount)
	} else {
		log.Println("\n[Step 2/4] No orphaned stocks to delete")
	}

	// Step 3: Re-sync all stocks from products
	log.Println("\n[Step 3/4] Re-syncing all stocks from products...")
	err = dao.DB_SyncStocksFromProducts()
	if err != nil {
		return fmt.Errorf("failed to sync stocks: %v", err)
	}
	log.Println("Successfully re-synced all stocks")

	// Step 4: Validate integrity
	log.Println("\n[Step 4/4] Validating stock integrity...")
	report, err := dao.DB_ValidateStockIntegrity()
	if err != nil {
		return fmt.Errorf("failed to validate integrity: %v", err)
	}

	log.Printf("\nValidation Report:")
	log.Printf("  Status: %v", report["status"])
	log.Printf("  Total Issues: %v", report["total_issues"])
	log.Printf("  Orphaned Stocks: %v", report["orphaned_stocks_count"])

	if issues, ok := report["issues"].([]string); ok && len(issues) > 0 {
		log.Println("\nIssues Found:")
		for i, issue := range issues {
			log.Printf("  %d. %s", i+1, issue)
		}
	}

	log.Println("\n=== Stock Integrity Fix Completed ===")

	if report["status"] == "ok" {
		log.Println("✅ All checks passed! Stock data is now consistent.")
		return nil
	} else {
		log.Println("⚠️  Some issues remain. Please review the report above.")
		return fmt.Errorf("stock integrity issues detected")
	}
}

// GetStockStatistics returns statistics about the stock collection
func GetStockStatistics() (map[string]interface{}, error) {
	ctx := context.Background()
	stocksCollection := dbConfigs.DATABASE.Collection("Stocks")

	stats := make(map[string]interface{})

	// Total stocks
	totalStocks, err := stocksCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats["total_stocks"] = totalStocks

	// Stocks with batches
	stocksWithBatches, err := stocksCollection.CountDocuments(ctx, bson.M{
		"batchId": bson.M{
			"$exists": true,
			"$nin":    []interface{}{"", nil},
		},
	})
	if err != nil {
		return nil, err
	}
	stats["stocks_with_batches"] = stocksWithBatches

	// Orphaned stocks
	orphanedStocks, err := stocksCollection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"batchId": nil},
			{"batchId": ""},
			{"batchId": bson.M{"$exists": false}},
		},
	})
	if err != nil {
		return nil, err
	}
	stats["orphaned_stocks"] = orphanedStocks

	// Unique products in stock
	productIds, err := stocksCollection.Distinct(ctx, "productId", bson.M{})
	if err != nil {
		return nil, err
	}
	stats["unique_products"] = len(productIds)

	return stats, nil
}
