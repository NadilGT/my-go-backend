package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
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
