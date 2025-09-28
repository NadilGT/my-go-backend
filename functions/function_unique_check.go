package functions

import (
	"context"
	"employee-crud/dbConfigs"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

func UniqueCheck(obj interface{}, collectionName string, fields []string) error {
	result := make(map[string]interface{})

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		if field.Kind() == reflect.String {
			result[fieldName] = field.String()
		}
	}
	var model map[string]interface{}
	filter := bson.M{
		"$or": make([]bson.M, 0, len(fields)),
	}
	for _, field := range fields {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{strings.ToLower(field): result[field]})
	}
	err := dbConfigs.DATABASE.Collection(collectionName).FindOne(context.Background(), filter).Decode(&model)

	if err != nil {
		return nil
	} else {
		return errors.New(fmt.Sprintf("%v unique", fields))
	}
}
