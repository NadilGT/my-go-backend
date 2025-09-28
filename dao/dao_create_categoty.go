package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_CreateCategory(object *dto.Category) error {

	_, err := dbConfigs.DATABASE.Collection("Categories").InsertOne(context.Background(), object)
	if err != nil {
		return err
	}
	return nil
}
