package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_CreateBrand(object *dto.Brand) error {
	_, err := dbConfigs.DATABASE.Collection("Brands").InsertOne(context.Background(), object)
	if err != nil {
		return err
	}
	return nil
}
