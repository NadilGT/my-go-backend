package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_AddBatchToProduct adds a new batch to an existing product
func DB_AddBatchToProduct(productId string, batch dto.Batch) error {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{"productId": productId}
	update := bson.M{
		"$push": bson.M{"batches": batch},
		"$inc":  bson.M{"stockQty": batch.StockQty},
		"$set":  bson.M{"updated_at": time.Now().UTC()},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DB_UpdateProductWithBatch updates a product's main fields and initializes batches array
func DB_UpdateProductWithBatch(product *dto.Product, initialBatch dto.Batch) error {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{"productId": product.ProductId}
	update := bson.M{
		"$set": bson.M{
			"batches":    []dto.Batch{initialBatch},
			"stockQty":   initialBatch.StockQty,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
