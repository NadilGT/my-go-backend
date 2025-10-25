package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_EditBatchStock edits the stock quantity of a specific batch
// Can increase or decrease the quantity
func DB_EditBatchStock(productId string, batchId string, newStockQty int) (*dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	if newStockQty < 0 {
		return nil, fmt.Errorf("stock quantity cannot be negative")
	}

	// Get the product
	var product dto.Product
	err := collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	// Find the batch and update it
	batchFound := false
	now := time.Now().UTC()

	for i := range product.Batches {
		if product.Batches[i].BatchId == batchId {
			product.Batches[i].StockQty = newStockQty
			product.Batches[i].UpdatedAt = now
			batchFound = true
			break
		}
	}

	if !batchFound {
		return nil, fmt.Errorf("batch not found: %s", batchId)
	}

	// If new quantity is 0, remove the batch
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

// DB_EditBatchDetails edits batch details including prices and expiry date
func DB_EditBatchDetails(productId string, batchId string, expiryDate *time.Time, costPrice float64, sellingPrice float64) (*dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	// Get the product
	var product dto.Product
	err := collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	// Find the batch and update it
	batchFound := false
	now := time.Now().UTC()

	for i := range product.Batches {
		if product.Batches[i].BatchId == batchId {
			if expiryDate != nil {
				product.Batches[i].ExpiryDate = expiryDate
			}
			if costPrice > 0 {
				product.Batches[i].CostPrice = costPrice
			}
			if sellingPrice > 0 {
				product.Batches[i].SellingPrice = sellingPrice
			}
			product.Batches[i].UpdatedAt = now
			batchFound = true
			break
		}
	}

	if !batchFound {
		return nil, fmt.Errorf("batch not found: %s", batchId)
	}

	// Update product in database
	filter := bson.M{"productId": productId, "deleted": false}
	update := bson.M{
		"$set": bson.M{
			"batches":    product.Batches,
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
