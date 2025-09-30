package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_CreateSupplier(object *dto.Supplier) error {

	_, err := dbConfigs.DATABASE.Collection("Suppliers").InsertOne(context.Background(), object)
	if err != nil {
		return err
	}
	return nil
}
