package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB_FindAllStocksCursorPaginated retrieves all stocks with cursor-based pagination
// This is optimized for large datasets (10000+ records)
// Stocks are sorted by productId first, then by batchId to group batches under same product
func DB_FindAllStocksCursorPaginated(limit int, cursor string) ([]dto.Stock, string, bool, error) {
	collection := dbConfigs.DATABASE.Collection("Stocks")
	ctx := context.Background()

	filter := bson.M{}

	// If cursor is provided, add it to filter for cursor-based pagination
	if cursor != "" {
		// Parse cursor (updated_at timestamp) and add to filter
		cursorTime, err := time.Parse("2006-01-02T15:04:05.000Z", cursor)
		if err != nil {
			// If parsing fails, try RFC3339 format
			cursorTime, err = time.Parse(time.RFC3339, cursor)
			if err != nil {
				// If still fails, return error
				return nil, "", false, err
			}
		}
		filter["updated_at"] = bson.M{"$lt": cursorTime}
	}

	// Set up find options
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	// Sort by productId first (ascending), then by updated_at (descending)
	// This groups batches of the same product together
	findOptions.SetSort(bson.D{
		{Key: "productId", Value: 1},
		{Key: "updated_at", Value: -1},
	})

	cursor_result, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, "", false, err
	}
	defer cursor_result.Close(ctx)

	var stocks []dto.Stock
	if err := cursor_result.All(ctx, &stocks); err != nil {
		return nil, "", false, err
	}

	// Calculate status for each stock
	for i := range stocks {
		stocks[i].CalculateStatus()
	}

	// Determine next cursor and if there are more pages
	var nextCursor string
	hasMore := false

	if len(stocks) > 0 {
		// Use the updated_at of the last stock as the next cursor
		lastStock := stocks[len(stocks)-1]
		nextCursor = lastStock.UpdatedAt.Format("2006-01-02T15:04:05.000Z")

		// Check if there are more stocks after this cursor
		checkFilter := bson.M{
			"updated_at": bson.M{"$lt": lastStock.UpdatedAt},
		}
		count, err := collection.CountDocuments(ctx, checkFilter)
		if err == nil && count > 0 {
			hasMore = true
		}
	}

	return stocks, nextCursor, hasMore, nil
}

// DB_GetStocksCount returns the total count of stocks
func DB_GetStocksCount() (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Stocks")
	ctx := context.Background()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}
