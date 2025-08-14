package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_FindEmployeeByID(id string) (*dto.Employee, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}
	var employee dto.Employee
	err = dbConfigs.UserCollection.FindOne(context.Background(), bson.M{"id": intID}).Decode(&employee)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return &employee, nil
}
