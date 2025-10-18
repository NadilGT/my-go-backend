package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateProductStock(productId string, quantitySold int) error {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"productId": productId}
	update := bson.M{
		"$inc": bson.M{"stockQty": -quantitySold},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
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
