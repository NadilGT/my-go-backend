package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_RemoveStockFromBatch removes/reduces stock from a specific batch
// If quantity to remove equals or exceeds batch stock, the batch is deleted
func DB_RemoveStockFromBatch(productId string, batchId string, quantityToRemove int) (*dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	if quantityToRemove <= 0 {
		return nil, fmt.Errorf("quantity to remove must be greater than 0")
	}

	// Get the product
	var product dto.Product
	err := collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	// Find the batch and reduce stock
	batchFound := false
	now := time.Now().UTC()

	for i := range product.Batches {
		if product.Batches[i].BatchId == batchId {
			if product.Batches[i].StockQty < quantityToRemove {
				return nil, fmt.Errorf("insufficient stock in batch %s: requested %d, available %d",
					batchId, quantityToRemove, product.Batches[i].StockQty)
			}
			product.Batches[i].StockQty -= quantityToRemove
			product.Batches[i].UpdatedAt = now
			batchFound = true
			break
		}
	}

	if !batchFound {
		return nil, fmt.Errorf("batch not found: %s", batchId)
	}

	// Remove batches with 0 stock
	var updatedBatches []dto.Batch
	for _, batch := range product.Batches {
		if batch.StockQty > 0 {
			updatedBatches = append(updatedBatches, batch)
		}
	}

	// Calculate new total stock
	totalStock := calculateTotalStock(updatedBatches)

	// Update product in database
	filter := bson.M{"productId": productId, "deleted": false}
	update := bson.M{
		"$set": bson.M{
			"batches":    updatedBatches,
			"stockQty":   totalStock,
			"updated_at": now,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Get updated product
	err = collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// DB_DeleteBatch completely deletes a batch from a product
func DB_DeleteBatch(productId string, batchId string) (*dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	// Get the product
	var product dto.Product
	err := collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	// Find and remove the batch
	batchFound := false
	var updatedBatches []dto.Batch

	for _, batch := range product.Batches {
		if batch.BatchId == batchId {
			batchFound = true
			// Skip this batch (delete it)
			continue
		}
		updatedBatches = append(updatedBatches, batch)
	}

	if !batchFound {
		return nil, fmt.Errorf("batch not found: %s", batchId)
	}

	// Calculate new total stock
	totalStock := calculateTotalStock(updatedBatches)
	now := time.Now().UTC()

	// Update product in database
	filter := bson.M{"productId": productId, "deleted": false}
	update := bson.M{
		"$set": bson.M{
			"batches":    updatedBatches,
			"stockQty":   totalStock,
			"updated_at": now,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Get updated product
	err = collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
