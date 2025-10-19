package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_RestoreProductByID(productId, categoryId, brandId, subCategoryId string) error {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{
		"productId": productId,
		"deleted":   true,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted":       false,
			"categoryId":    categoryId,
			"brandId":       brandId,
			"subCategoryId": subCategoryId,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("Product not found or already active")
	}

	// Fetch the restored product and sync to Stocks collection
	var product dto.Product
	err = collection.FindOne(ctx, bson.M{"productId": productId}).Decode(&product)
	if err == nil {
		// Sync to stocks (ignore error to not fail the restore operation)
		DB_SyncSingleProductStock(&product)
	}

	return nil
}
