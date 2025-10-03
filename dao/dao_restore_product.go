package dao

import (
	"context"
	"employee-crud/dbConfigs"
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

	return nil
}
