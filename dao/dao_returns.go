package dao

import (
	"context"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ReturnsCollection *mongo.Collection

func InitReturnsCollection(db *mongo.Database) {
	ReturnsCollection = db.Collection("returns")
}

func InsertReturn(ctx context.Context, ret dto.ReturnDTO) error {
	_, err := ReturnsCollection.InsertOne(ctx, ret)
	return err
}

func GetAllReturns(ctx context.Context) ([]dto.ReturnDTO, error) {
	cursor, err := ReturnsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []dto.ReturnDTO
	for cursor.Next(ctx) {
		var r dto.ReturnDTO
		if err := cursor.Decode(&r); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func GetReturnByID(ctx context.Context, id string) (*dto.ReturnDTO, error) {
	var r dto.ReturnDTO
	err := ReturnsCollection.FindOne(ctx, bson.M{"id": id}).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
