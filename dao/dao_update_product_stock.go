package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateProductStock(productId string, quantitySold int) error {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"productId": productId, "deleted": false}
	update := bson.M{
		"$inc": bson.M{"stockQty": -quantitySold},
		"$set": bson.M{"updated_at": time.Now()},
	}

	// Use FindOneAndUpdate to get the updated document in a single operation
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedProduct dto.Product
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedProduct)
	if err != nil {
		return err
	}

	// Sync the updated stock to Stocks collection
	// This ensures the Stocks collection is immediately updated when a sale is made
	if err := DB_SyncSingleProductStock(&updatedProduct); err != nil {
		// Return the error to ensure sync failures are caught
		return err
	}

	return nil
}

func GetProductByProductId(productId string) (*dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product dto.Product
	err := collection.FindOne(ctx, bson.M{"productId": productId, "deleted": false}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
