package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DB_FindAllProducts() ([]dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

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

func DB_FindAllProductsPaginated(page, limit int) ([]dto.Product, int64, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	// Calculate skip value
	skip := (page - 1) * limit

	// Count total documents matching the filter
	// For better performance with large datasets, consider caching this count
	// or using estimatedDocumentCount() if exact count isn't critical
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Set up find options with pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	// Sort by created_at descending for consistent results
	// For better performance with large datasets, ensure you have an index on created_at
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []dto.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Cursor-based pagination for maximum performance with large datasets
func DB_FindAllProductsCursorPaginated(limit int, cursor string) ([]dto.Product, string, bool, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{"deleted": false}

	// If cursor is provided, add it to filter for cursor-based pagination
	if cursor != "" {
		// Parse cursor (created_at timestamp) and add to filter
		filter["created_at"] = bson.M{"$lt": cursor}
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
