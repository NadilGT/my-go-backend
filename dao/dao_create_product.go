package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_CreateProduct(object *dto.Product) error {

	_, err := dbConfigs.DATABASE.Collection("Products").InsertOne(context.Background(), object)
	if err != nil {
		return err
	}
	return nil
}
