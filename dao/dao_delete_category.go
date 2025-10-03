package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DB_DeletecategoryByID(categoryId string) error {
	ctx := context.Background()

	categories := dbConfigs.DATABASE.Collection("Categories")
	products := dbConfigs.DATABASE.Collection("Products")

	filter := bson.M{
		"categoryId": categoryId,
		"deleted":    false,
	}

	update := bson.M{
		"$set": bson.M{"deleted": true},
	}

	result, err := categories.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("Specified category not found or already deleted")
	}

	productFilter := bson.M{
		"categoryId": categoryId,
		"deleted":    false,
	}

	productUpdate := bson.M{
		"$set": bson.M{"deleted": true},
	}

	_, err = products.UpdateMany(ctx, productFilter, productUpdate)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	return nil
}
