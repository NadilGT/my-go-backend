package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB_SyncStocksFromProducts syncs all product stocks to the Stocks collection
// This function handles large datasets efficiently by using batch processing
func DB_SyncStocksFromProducts() error {
	productsCollection := dbConfigs.DATABASE.Collection("Products")
	stocksCollection := dbConfigs.DATABASE.Collection("Stocks")
	ctx := context.Background()

	// Filter only non-deleted products
	filter := bson.M{"deleted": false}

	// Find all products in batches for memory efficiency
	batchSize := 500 // Process 500 products at a time
	skip := 0

	for {
		// Set up find options with batch processing
		findOptions := options.Find()
		findOptions.SetSkip(int64(skip))
		findOptions.SetLimit(int64(batchSize))

		cursor, err := productsCollection.Find(ctx, filter, findOptions)
		if err != nil {
			return err
		}

		var products []dto.Product
		if err := cursor.All(ctx, &products); err != nil {
			cursor.Close(ctx)
			return err
		}
		cursor.Close(ctx)

		// If no more products, break the loop
		if len(products) == 0 {
			break
		}

		// Prepare bulk write operations for better performance
		var operations []mongo.WriteModel
		currentTime := time.Now()

		for _, product := range products {
			// Only sync products that have batches
			// Skip products without batches to prevent null batchId entries
			if len(product.Batches) > 0 {
				for _, batch := range product.Batches {
					// Validate batch has an ID
					if batch.BatchId == "" {
						continue // Skip batches with empty IDs
					}

					filter := bson.M{
						"productId": product.ProductId,
						"batchId":   batch.BatchId,
					}
					update := bson.M{
						"$set": bson.M{
							"productId":   product.ProductId,
							"batchId":     batch.BatchId,
							"name":        product.Name,
							"stockQty":    batch.StockQty,
							"expiry_date": batch.ExpiryDate,
							"updated_at":  currentTime,
						},
						"$setOnInsert": bson.M{
							"created_at": currentTime,
						},
					}

					operation := mongo.NewUpdateOneModel()
					operation.SetFilter(filter)
					operation.SetUpdate(update)
					operation.SetUpsert(true)

					operations = append(operations, operation)
				}
			}
		}

		// Execute bulk write
		if len(operations) > 0 {
			_, err := stocksCollection.BulkWrite(ctx, operations)
			if err != nil {
				return err
			}
		}

		// Move to next batch
		skip += batchSize
	}

	return nil
}

// DB_SyncSingleProductStock syncs a single product's stock to the Stocks collection
// Use this when a product is created or updated
func DB_SyncSingleProductStock(product *dto.Product) error {
	stocksCollection := dbConfigs.DATABASE.Collection("Stocks")
	ctx := context.Background()

	currentTime := time.Now()

	// If product has batches, sync each batch as a separate stock entry
	if len(product.Batches) > 0 {
		// First, delete ALL existing stock entries for this product (including orphaned null batches)
		// Then recreate them from the current batches - this is cleaner and ensures no orphaned data
		stocksCollection.DeleteMany(ctx, bson.M{"productId": product.ProductId})

		// Sync each batch
		for _, batch := range product.Batches {
			// Use productId + batchId as unique identifier
			filter := bson.M{
				"productId": product.ProductId,
				"batchId":   batch.BatchId,
			}
			update := bson.M{
				"$set": bson.M{
					"productId":   product.ProductId,
					"batchId":     batch.BatchId,
					"name":        product.Name,
					"stockQty":    batch.StockQty,
					"expiry_date": batch.ExpiryDate,
					"updated_at":  currentTime,
				},
				"$setOnInsert": bson.M{
					"created_at": currentTime,
				},
			}

			opts := options.Update().SetUpsert(true)
			_, err := stocksCollection.UpdateOne(ctx, filter, update, opts)
			if err != nil {
				return err
			}
		}
	} else {
		// Products without batches: Skip syncing to avoid orphaned stock entries
		// Products should use the batch system - this is legacy code path
		// Delete any existing stock entries for this product
		stocksCollection.DeleteMany(ctx, bson.M{"productId": product.ProductId})
		// Note: Products without batches won't appear in stock listing
		// This encourages migration to batch-based inventory system
	}

	return nil
}
