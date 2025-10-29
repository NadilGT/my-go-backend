package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DB_FindProductsByBarcodeCursorPaginated(barcode string, limit int, cursor string) ([]dto.Product, string, bool, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{
		"barcode": bson.M{"$regex": barcode, "$options": "i"},
		"deleted": false,
	}
	
	if cursor != "" {
		cursorTime, err := time.Parse("2006-01-02T15:04:05.000Z", cursor)
		if err != nil {
			cursorTime, err = time.Parse(time.RFC3339, cursor)
			if err != nil {
				return nil, "", false, err
			}
		}
		filter["created_at"] = bson.M{"$lt": cursorTime}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursorResult, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, "", false, err
	}
	defer cursorResult.Close(ctx)

	var products []dto.Product
	if err := cursorResult.All(ctx, &products); err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasMore := false

	if len(products) > 0 {
		lastProduct := products[len(products)-1]
		nextCursor = lastProduct.CreatedAt.Format("2006-01-02T15:04:05.000Z")

		checkFilter := bson.M{
			"barcode": bson.M{"$regex": barcode, "$options": "i"},
			"deleted": false,
			"created_at": bson.M{
				"$lt": lastProduct.CreatedAt,
			},
		}

		count, err := collection.CountDocuments(ctx, checkFilter)
		if err == nil && count > 0 {
			hasMore = true
		}
	}

	return products, nextCursor, hasMore, nil
}
