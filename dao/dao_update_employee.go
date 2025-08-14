package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_UpdateEmployee(employee *dto.Employee) error {
	_, err := dbConfigs.UserCollection.UpdateOne(context.Background(), bson.M{"id": employee.ID}, bson.M{"$set": employee})
	if err != nil {
		return err
	}
	return nil
}
