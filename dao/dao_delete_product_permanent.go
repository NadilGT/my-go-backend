package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_PermanentDeleteProductByID(productId string) error {
	collection := dbConfigs.DATABASE.Collection("Products")

	filter := bson.M{"productId": productId}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("Specified productId not found")
	}

	return nil
}
