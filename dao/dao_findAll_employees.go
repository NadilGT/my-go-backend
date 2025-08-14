package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_FindAllEmployees() (*[]dto.Employee, error) {
	var employees []dto.Employee

	filter := bson.D{{Key: "deleted", Value: false}}
	results, err := dbConfigs.UserCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	for results.Next(context.Background()) {
		var employee dto.Employee
		if err := results.Decode(&employee); err != nil {
			return nil, errors.New("error decoding employee")
		}
		employees = append(employees, employee)
	}
	return &employees, nil
}
