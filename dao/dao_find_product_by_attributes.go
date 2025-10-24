package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
)

// DB_FindProductByAttributes finds a product with matching category, brand, and subcategory
func DB_FindProductByAttributes(categoryId, brandId, subCategoryId string) (*dto.Product, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	filter := bson.M{
		"categoryId":    categoryId,
		"brandId":       brandId,
		"subCategoryId": subCategoryId,
		"deleted":       false,
	}

	var product dto.Product
	err := collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
