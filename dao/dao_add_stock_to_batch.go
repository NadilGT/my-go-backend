package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_AddStockToProduct adds stock to an existing product
// If expiry date matches existing batch, adds to that batch
// If expiry date is different, creates a new batch
func DB_AddStockToProduct(productId string, stockQty int, expiryDate *time.Time, costPrice float64, sellingPrice float64) (*dto.Product, string, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	// Get the product
	var product dto.Product
	err := collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return nil, "", fmt.Errorf("product not found: %v", err)
	}

	now := time.Now().UTC()

	// Check if product has batches
	if len(product.Batches) > 0 {
		// Look for a batch with matching expiry date
		batchFound := false
		for i := range product.Batches {
			if datesMatch(product.Batches[i].ExpiryDate, expiryDate) {
				// Found matching batch - add stock to it
				product.Batches[i].StockQty += stockQty
				product.Batches[i].UpdatedAt = now
				// Update prices if provided and different
				if costPrice > 0 {
					product.Batches[i].CostPrice = costPrice
				}
				if sellingPrice > 0 {
					product.Batches[i].SellingPrice = sellingPrice
				}
				batchFound = true

				// Update product in database
				filter := bson.M{"productId": productId, "deleted": false}
				update := bson.M{
					"$set": bson.M{
						"batches":    product.Batches,
						"stockQty":   calculateTotalStock(product.Batches),
						"updated_at": now,
					},
				}

				_, err := collection.UpdateOne(ctx, filter, update)
				if err != nil {
					return nil, "", err
				}

				// Get updated product
				err = collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
				if err != nil {
					return nil, "", err
				}

				return &product, product.Batches[i].BatchId, nil
			}
		}

		if !batchFound {
			// No matching expiry date - create new batch
			batchId, err := GenerateId(ctx, "Batches", "BATCH")
			if err != nil {
				return nil, "", err
			}

			newBatch := dto.Batch{
				BatchId:      batchId,
				StockQty:     stockQty,
				ExpiryDate:   expiryDate,
				CostPrice:    costPrice,
				SellingPrice: sellingPrice,
				CreatedAt:    now,
				UpdatedAt:    now,
			}

			// Add new batch
			if err := DB_AddBatchToProduct(productId, newBatch); err != nil {
				return nil, "", err
			}

			// Get updated product
			err = collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
			if err != nil {
				return nil, "", err
			}

			return &product, batchId, nil
		}
	} else {
		// Product has no batches - create first batch
		batchId, err := GenerateId(ctx, "Batches", "BATCH")
		if err != nil {
			return nil, "", err
		}

		newBatch := dto.Batch{
			BatchId:      batchId,
			StockQty:     stockQty,
			ExpiryDate:   expiryDate,
			CostPrice:    costPrice,
			SellingPrice: sellingPrice,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// Initialize batches array with first batch
		if err := DB_UpdateProductWithBatch(&product, newBatch); err != nil {
			return nil, "", err
		}

		// Get updated product
		err = collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
		if err != nil {
			return nil, "", err
		}

		return &product, batchId, nil
	}

	return &product, "", nil
}

// Helper function to check if two dates match (ignoring time)
func datesMatch(date1 *time.Time, date2 *time.Time) bool {
	if date1 == nil && date2 == nil {
		return true
	}
	if date1 == nil || date2 == nil {
		return false
	}
	return date1.Format("2006-01-02") == date2.Format("2006-01-02")
}

// Helper function to calculate total stock from batches
func calculateTotalStock(batches []dto.Batch) int {
	total := 0
	for _, batch := range batches {
		total += batch.StockQty
	}
	return total
}
