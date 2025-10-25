package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductWithStockInfo represents a product with its stock information
type ProductWithStockInfo struct {
	ProductId  string     `json:"productId"`
	Name       string     `json:"name"`
	StockQty   int        `json:"stockQty"`
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`
	BatchId    string     `json:"batchId,omitempty"`
	HasBatches bool       `json:"hasBatches"`
	BatchCount int        `json:"batchCount"`
	Status     string     `json:"status,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// CalculateStatus calculates the stock status for a product
func (p *ProductWithStockInfo) CalculateStatus() {
	if p.StockQty < 10 {
		p.Status = "Low"
	} else if p.StockQty >= 10 && p.StockQty < 25 {
		p.Status = "Average"
	} else {
		p.Status = "Good"
	}
}

// DB_FindAllProductsWithStock retrieves all products with their stock information
// This includes products with and without batches
func DB_FindAllProductsWithStock(limit int, cursor string) ([]ProductWithStockInfo, string, bool, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	// If cursor is provided, add it to filter for cursor-based pagination
	if cursor != "" {
		// Parse cursor as timestamp
		cursorTime, err := time.Parse(time.RFC3339Nano, cursor)
		if err != nil {
			// Try alternative format
			cursorTime, err = time.Parse(time.RFC3339, cursor)
			if err != nil {
				return nil, "", false, err
			}
		}
		filter["updated_at"] = bson.M{"$lt": cursorTime}
	}

	// Set up find options
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{
		{Key: "updated_at", Value: -1},
		{Key: "_id", Value: -1},
	})

	cursor_result, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, "", false, err
	}
	defer cursor_result.Close(ctx)

	var products []dto.Product
	if err := cursor_result.All(ctx, &products); err != nil {
		return nil, "", false, err
	}

	// Convert products to ProductWithStockInfo
	var productsWithStock []ProductWithStockInfo
	for _, product := range products {
		if len(product.Batches) > 0 {
			// Product has batches - create entry for each batch
			for _, batch := range product.Batches {
				stockInfo := ProductWithStockInfo{
					ProductId:  product.ProductId,
					Name:       product.Name,
					StockQty:   batch.StockQty,
					ExpiryDate: batch.ExpiryDate,
					BatchId:    batch.BatchId,
					HasBatches: true,
					BatchCount: len(product.Batches),
					CreatedAt:  product.CreatedAt,
					UpdatedAt:  product.UpdatedAt,
				}
				stockInfo.CalculateStatus()
				productsWithStock = append(productsWithStock, stockInfo)
			}
		} else {
			// Product has no batches - create single entry
			stockInfo := ProductWithStockInfo{
				ProductId:  product.ProductId,
				Name:       product.Name,
				StockQty:   product.StockQty,
				ExpiryDate: product.ExpiryDate,
				BatchId:    "", // No batch
				HasBatches: false,
				BatchCount: 0,
				CreatedAt:  product.CreatedAt,
				UpdatedAt:  product.UpdatedAt,
			}
			stockInfo.CalculateStatus()
			productsWithStock = append(productsWithStock, stockInfo)
		}
	}

	// Determine next cursor and if there are more pages
	var nextCursor string
	hasMore := false

	if len(products) > 0 {
		lastProduct := products[len(products)-1]
		// Use only UpdatedAt for cursor since Product doesn't expose _id
		nextCursor = lastProduct.UpdatedAt.Format(time.RFC3339Nano)

		// Check if there are more products
		checkFilter := bson.M{
			"deleted":    false,
			"updated_at": bson.M{"$lt": lastProduct.UpdatedAt},
		}

		count, err := collection.CountDocuments(ctx, checkFilter)
		if err == nil && count > 0 {
			hasMore = true
		}
	}

	return productsWithStock, nextCursor, hasMore, nil
}

// DB_GetProductsWithStockCount returns the total count of non-deleted products
func DB_GetProductsWithStockCount() (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	count, err := collection.CountDocuments(ctx, bson.M{"deleted": false})
	if err != nil {
		return 0, err
	}

	return count, nil
}
