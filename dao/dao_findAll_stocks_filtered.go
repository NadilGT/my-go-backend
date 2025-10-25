package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB_FindAllStocksFilteredCursorPaginated retrieves ALL batches from products filtered by total stockQty
// This filters PRODUCTS by their total stock status, then returns ALL their batches
// Parameters:
//   - limit: number of records per page
//   - cursor: cursor for pagination (updated_at timestamp)
//   - minQty: minimum stock quantity (inclusive) - applied to PRODUCT total
//   - maxQty: maximum stock quantity (inclusive, use -1 for no upper limit) - applied to PRODUCT total
func DB_FindAllStocksFilteredCursorPaginated(limit int, cursor string, minQty int, maxQty int) ([]dto.Stock, string, bool, error) {
	productsCollection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	// Build filter to find products by their TOTAL stockQty
	filter := bson.M{"deleted": false}

	// Add quantity filter on product's total stockQty
	if maxQty == -1 {
		// No upper limit (e.g., for "good" status)
		filter["stockQty"] = bson.M{"$gte": minQty}
	} else {
		// Range filter
		filter["stockQty"] = bson.M{
			"$gte": minQty,
			"$lte": maxQty,
		}
	}

	// If cursor is provided, add it to filter for cursor-based pagination
	if cursor != "" {
		cursorTime, err := time.Parse("2006-01-02T15:04:05.000Z", cursor)
		if err != nil {
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
	})

	cursor_result, err := productsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, "", false, err
	}
	defer cursor_result.Close(ctx)

	var products []dto.Product
	if err := cursor_result.All(ctx, &products); err != nil {
		return nil, "", false, err
	}

	// Convert products to stocks (include ALL batches from each product)
	var stocks []dto.Stock
	for _, product := range products {
		if len(product.Batches) > 0 {
			// Product has batches - create stock entry for each batch
			for _, batch := range product.Batches {
				stock := dto.Stock{
					ProductId:  product.ProductId,
					Name:       product.Name,
					BatchId:    batch.BatchId,
					StockQty:   batch.StockQty,
					ExpiryDate: batch.ExpiryDate,
					CreatedAt:  batch.CreatedAt,
					UpdatedAt:  batch.UpdatedAt,
				}
				stock.CalculateStatus()
				stocks = append(stocks, stock)
			}
		} else {
			// Product has no batches - create single stock entry
			stock := dto.Stock{
				ProductId:  product.ProductId,
				Name:       product.Name,
				BatchId:    "",
				StockQty:   product.StockQty,
				ExpiryDate: product.ExpiryDate,
				CreatedAt:  product.CreatedAt,
				UpdatedAt:  product.UpdatedAt,
			}
			stock.CalculateStatus()
			stocks = append(stocks, stock)
		}
	}

	// Determine next cursor and if there are more pages
	var nextCursor string
	hasMore := false

	if len(products) > 0 {
		lastProduct := products[len(products)-1]
		nextCursor = lastProduct.UpdatedAt.Format("2006-01-02T15:04:05.000Z")

		// Check if there are more products after this cursor with the same filter
		checkFilter := bson.M{
			"deleted":    false,
			"updated_at": bson.M{"$lt": lastProduct.UpdatedAt},
		}

		// Apply the same quantity filter
		if maxQty == -1 {
			checkFilter["stockQty"] = bson.M{"$gte": minQty}
		} else {
			checkFilter["stockQty"] = bson.M{
				"$gte": minQty,
				"$lte": maxQty,
			}
		}

		count, err := productsCollection.CountDocuments(ctx, checkFilter)
		if err == nil && count > 0 {
			hasMore = true
		}
	}

	return stocks, nextCursor, hasMore, nil
}

// DB_GetStocksCountFiltered returns the count of products filtered by quantity range
// This counts PRODUCTS, not individual batches
func DB_GetStocksCountFiltered(minQty int, maxQty int) (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	// Build filter based on product's total stockQty
	filter := bson.M{"deleted": false}
	if maxQty == -1 {
		// No upper limit
		filter["stockQty"] = bson.M{"$gte": minQty}
	} else {
		// Range filter
		filter["stockQty"] = bson.M{
			"$gte": minQty,
			"$lte": maxQty,
		}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
