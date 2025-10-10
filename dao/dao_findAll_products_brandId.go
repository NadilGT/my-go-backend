package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DB_FindProductsByBrand(brandId string) ([]dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{
		"brandId": brandId,
		"deleted": false,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []dto.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// Cursor-based pagination for products by brand
func DB_FindProductsByBrandCursorPaginated(brandId string, limit int, cursor string) ([]dto.Product, string, bool, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{
		"brandId": brandId,
		"deleted": false,
	}

	// If cursor is provided, add it to filter for cursor-based pagination
	if cursor != "" {
		// Parse cursor (created_at timestamp) and add to filter
		cursorTime, err := time.Parse("2006-01-02T15:04:05.000Z", cursor)
		if err != nil {
			// If parsing fails, try RFC3339 format
			cursorTime, err = time.Parse(time.RFC3339, cursor)
			if err != nil {
				// If still fails, return error
				return nil, "", false, err
			}
		}
		filter["created_at"] = bson.M{"$lt": cursorTime}
	}

	// Set up find options
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor_result, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, "", false, err
	}
	defer cursor_result.Close(ctx)

	var products []dto.Product
	if err := cursor_result.All(ctx, &products); err != nil {
		return nil, "", false, err
	}

	// Determine next cursor and if there are more pages
	var nextCursor string
	hasMore := false

	if len(products) > 0 {
		// Use the created_at of the last product as the next cursor
		lastProduct := products[len(products)-1]
		nextCursor = lastProduct.CreatedAt.Format("2006-01-02T15:04:05.000Z")

		// Check if there are more products after this cursor
		checkFilter := bson.M{
			"brandId":    brandId,
			"deleted":    false,
			"created_at": bson.M{"$lt": lastProduct.CreatedAt},
		}
		count, err := collection.CountDocuments(ctx, checkFilter)
		if err == nil && count > 0 {
			hasMore = true
		}
	}

	return products, nextCursor, hasMore, nil
}
