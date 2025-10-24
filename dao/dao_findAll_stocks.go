package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"encoding/base64"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CursorData holds the compound cursor data for pagination
type CursorData struct {
	UpdatedAt time.Time `json:"updated_at"`
	ID        string    `json:"id"`
}

// encodeCursor creates a base64-encoded compound cursor from timestamp and ID
func encodeCursor(updatedAt time.Time, id primitive.ObjectID) string {
	cursorData := CursorData{
		UpdatedAt: updatedAt,
		ID:        id.Hex(),
	}
	jsonData, err := json.Marshal(cursorData)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(jsonData)
}

// decodeCursor decodes a base64-encoded compound cursor
func decodeCursor(cursor string) (*CursorData, error) {
	// Try to decode as compound cursor (new format)
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err == nil {
		var cursorData CursorData
		if err := json.Unmarshal(decoded, &cursorData); err == nil {
			return &cursorData, nil
		}
	}

	// Fallback: Try to parse as old format (plain timestamp) for backward compatibility
	cursorTime, err := time.Parse("2006-01-02T15:04:05.000Z", cursor)
	if err != nil {
		cursorTime, err = time.Parse(time.RFC3339, cursor)
		if err != nil {
			return nil, err
		}
	}

	// Return cursor data with empty ID (old format compatibility)
	return &CursorData{
		UpdatedAt: cursorTime,
		ID:        "",
	}, nil
}

// DB_FindAllStocksCursorPaginated retrieves all stocks with cursor-based pagination
// This is optimized for large datasets (10000+ records)
// Uses compound cursor (updated_at + id) to handle duplicate timestamps
func DB_FindAllStocksCursorPaginated(limit int, cursor string) ([]dto.Stock, string, bool, error) {
	collection := dbConfigs.DATABASE.Collection("Stocks")
	ctx := context.Background()

	filter := bson.M{}

	// If cursor is provided, add it to filter for cursor-based pagination
	if cursor != "" {
		cursorData, err := decodeCursor(cursor)
		if err != nil {
			return nil, "", false, err
		}

		// Build compound cursor filter to handle duplicate timestamps
		// Query: WHERE (updated_at < cursor_time) OR (updated_at = cursor_time AND _id < cursor_id)
		if cursorData.ID != "" {
			// New format with compound cursor
			cursorObjID, err := primitive.ObjectIDFromHex(cursorData.ID)
			if err != nil {
				return nil, "", false, err
			}

			filter["$or"] = []bson.M{
				{"updated_at": bson.M{"$lt": cursorData.UpdatedAt}},
				{
					"updated_at": cursorData.UpdatedAt,
					"_id":        bson.M{"$lt": cursorObjID},
				},
			}
		} else {
			// Old format compatibility (only timestamp)
			filter["updated_at"] = bson.M{"$lt": cursorData.UpdatedAt}
		}
	}

	// Set up find options
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	// Sort by updated_at (descending), then by _id (descending) for consistent cursor pagination
	// Note: Sorting by productId was removed as it conflicts with cursor-based pagination
	findOptions.SetSort(bson.D{
		{Key: "updated_at", Value: -1},
		{Key: "_id", Value: -1},
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
		// Create compound cursor using updated_at and _id of the last stock
		lastStock := stocks[len(stocks)-1]
		nextCursor = encodeCursor(lastStock.UpdatedAt, lastStock.ID)

		// Check if there are more stocks after this cursor
		checkFilter := bson.M{
			"$or": []bson.M{
				{"updated_at": bson.M{"$lt": lastStock.UpdatedAt}},
				{
					"updated_at": lastStock.UpdatedAt,
					"_id":        bson.M{"$lt": lastStock.ID},
				},
			},
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
