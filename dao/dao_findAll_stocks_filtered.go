package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB_FindAllStocksFilteredCursorPaginated retrieves stocks filtered by quantity with cursor-based pagination
// This allows filtering stocks by status (low, average, good) across all pages
// Parameters:
//   - limit: number of records per page
//   - cursor: cursor for pagination (updated_at timestamp)
//   - minQty: minimum stock quantity (inclusive)
//   - maxQty: maximum stock quantity (inclusive, use -1 for no upper limit)
func DB_FindAllStocksFilteredCursorPaginated(limit int, cursor string, minQty int, maxQty int) ([]dto.Stock, string, bool, error) {
	collection := dbConfigs.DATABASE.Collection("Stocks")
	ctx := context.Background()

	// Build filter based on quantity range
	filter := bson.M{}

	// Add quantity filter
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

		// Check if there are more stocks after this cursor with the same filter
		checkFilter := bson.M{
			"updated_at": bson.M{"$lt": lastStock.UpdatedAt},
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

		count, err := collection.CountDocuments(ctx, checkFilter)
		if err == nil && count > 0 {
			hasMore = true
		}
	}

	return stocks, nextCursor, hasMore, nil
}

// DB_GetStocksCountFiltered returns the count of stocks filtered by quantity range
func DB_GetStocksCountFiltered(minQty int, maxQty int) (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Stocks")
	ctx := context.Background()

	// Build filter based on quantity range
	filter := bson.M{}
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
