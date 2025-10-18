package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func FindAllSales(limit int64, offset int64) ([]dto.Sale, error) {
	collection := dbConfigs.DATABASE.Collection("Sales")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := bson.M{}

	cursor, err := collection.Find(ctx, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []dto.Sale
	if err = cursor.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, nil
}
