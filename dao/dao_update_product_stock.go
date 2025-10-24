package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateProductStock(productId string, quantitySold int) error {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, get the product to check if it has batches
	var product dto.Product
	err := collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return err
	}

	// If product has batches, deduct from batches using FEFO (First Expired First Out)
	if len(product.Batches) > 0 {
		// Sort batches by expiry date (earliest first) for FEFO
		sort.Slice(product.Batches, func(i, j int) bool {
			// Handle nil expiry dates (put them at the end)
			if product.Batches[i].ExpiryDate == nil {
				return false
			}
			if product.Batches[j].ExpiryDate == nil {
				return true
			}
			return product.Batches[i].ExpiryDate.Before(*product.Batches[j].ExpiryDate)
		})

		remainingQty := quantitySold

		// Process batches in order
		var updatedBatches []dto.Batch
		for _, batch := range product.Batches {
			if remainingQty <= 0 {
				updatedBatches = append(updatedBatches, batch)
				continue
			}

			if batch.StockQty >= remainingQty {
				// This batch has enough stock
				batch.StockQty -= remainingQty
				batch.UpdatedAt = time.Now().UTC()
				remainingQty = 0
				if batch.StockQty > 0 {
					updatedBatches = append(updatedBatches, batch)
				}
			} else {
				// This batch doesn't have enough stock, use all of it
				remainingQty -= batch.StockQty
				// Don't add this batch to updated list (it's empty)
			}
		}

		if remainingQty > 0 {
			return fmt.Errorf("insufficient stock: requested %d, available %d", quantitySold, quantitySold-remainingQty)
		}

		// Calculate new total stock qty
		var totalStockQty int
		for _, batch := range updatedBatches {
			totalStockQty += batch.StockQty
		}

		// Update product with new batches and stock quantity
		filter := bson.M{"productId": productId, "deleted": false}
		update := bson.M{
			"$set": bson.M{
				"batches":    updatedBatches,
				"stockQty":   totalStockQty,
				"updated_at": time.Now().UTC(),
			},
		}

		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		// Fetch updated product for sync
		err = collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
		if err != nil {
			return err
		}

		// Sync the updated stock to Stocks collection
		if err := DB_SyncSingleProductStock(&product); err != nil {
			return err
		}
	} else {
		// Legacy: product without batches
		filter := bson.M{"productId": productId, "deleted": false}
		update := bson.M{
			"$inc": bson.M{"stockQty": -quantitySold},
			"$set": bson.M{"updated_at": time.Now()},
		}

		// Use FindOneAndUpdate to get the updated document
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		var updatedProduct dto.Product
		err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedProduct)
		if err != nil {
			return err
		}

		// Sync the updated stock to Stocks collection
		if err := DB_SyncSingleProductStock(&updatedProduct); err != nil {
			return err
		}
	}

	return nil
}

func GetProductByProductId(productId string) (*dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product dto.Product
	err := collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
