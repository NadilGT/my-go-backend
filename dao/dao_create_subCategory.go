package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_CreateSubCategory(object *dto.SubCategory) error {

	_, err := dbConfigs.DATABASE.Collection("SubCategories").InsertOne(context.Background(), object)
	if err != nil {
		return err
	}
	return nil
}
