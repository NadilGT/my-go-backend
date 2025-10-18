package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"
)

func CreateSale(sale *dto.Sale) error {
	collection := dbConfigs.DATABASE.Collection("Sales")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, sale)
	return err
}
