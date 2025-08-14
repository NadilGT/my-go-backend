package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_CreateEmployee(user dto.Employee) error {

	filter := bson.M{"email": user.Email}
	count, err := dbConfigs.UserCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already exists")
	}

	_, err = dbConfigs.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}
