package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_SoftDeleteEmployeeByID(id string) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"deleted": true}}
	_, err := dbConfigs.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("soft delete failed: %v", err)
	}
	return nil
}

func DB_HardDeleteEmployeeByID(id string) error {
	filter := bson.M{"id": id}
	_, err := dbConfigs.UserCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("hard delete failed: %v", err)
	}
	return nil
}
