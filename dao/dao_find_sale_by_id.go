package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func FindSaleBySaleId(saleId string) (*dto.Sale, error) {
	collection := dbConfigs.DATABASE.Collection("Sales")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sale dto.Sale
	err := collection.FindOne(ctx, bson.M{"saleId": saleId}).Decode(&sale)
	if err != nil {
		return nil, err
	}

	return &sale, nil
}
