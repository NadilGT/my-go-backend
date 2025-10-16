package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
)

func DB_CreateGRN(object *dto.GRN) error {
	_, err := dbConfigs.DATABASE.Collection("GRNs").InsertOne(context.Background(), object)
	if err != nil {
		return err
	}
	return nil
}
