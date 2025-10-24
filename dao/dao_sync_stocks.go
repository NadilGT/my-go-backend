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
			// If product has batches, create separate stock entries for each batch
			if len(product.Batches) > 0 {
				for _, batch := range product.Batches {
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
			} else {
				// Legacy support: products without batches
				filter := bson.M{
					"productId": product.ProductId,
					"batchId":   bson.M{"$exists": false},
				}
				update := bson.M{
					"$set": bson.M{
						"productId":   product.ProductId,
						"name":        product.Name,
						"stockQty":    product.StockQty,
						"expiry_date": product.ExpiryDate,
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
		// First, get all existing stock entries for this product to clean up removed batches
		existingStocks, err := stocksCollection.Find(ctx, bson.M{"productId": product.ProductId})
		if err == nil {
			var stocks []dto.Stock
			if err := existingStocks.All(ctx, &stocks); err == nil {
				// Create a map of current batch IDs
				currentBatchIds := make(map[string]bool)
				for _, batch := range product.Batches {
					currentBatchIds[batch.BatchId] = true
				}

				// Delete stock entries for batches that no longer exist
				for _, stock := range stocks {
					if stock.BatchId != "" && !currentBatchIds[stock.BatchId] {
						stocksCollection.DeleteOne(ctx, bson.M{"_id": stock.ID})
					}
				}
			}
			existingStocks.Close(ctx)
		}

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
		// Legacy support: products without batches
		filter := bson.M{
			"productId": product.ProductId,
			"batchId":   bson.M{"$exists": false},
		}
		update := bson.M{
			"$set": bson.M{
				"productId":   product.ProductId,
				"name":        product.Name,
				"stockQty":    product.StockQty,
				"expiry_date": product.ExpiryDate,
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

	return nil
}
