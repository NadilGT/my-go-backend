package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DB_UpdateProduct(ctx context.Context, product *dto.Product) error {
	collection := dbConfigs.DATABASE.Collection("Products")

	filter := bson.M{"productId": product.ProductId}

	update := bson.M{
		"$set": bson.M{
			"name":          product.Name,
			"barcode":       product.Barcode,
			"categoryId":    product.CategoryID,
			"brandId":       product.BrandID,
			"subCategoryId": product.SubCategoryID,
			"costPrice":     product.CostPrice,
			"sellingPrice":  product.SellingPrice,
			"stockQty":      product.StockQty,
			"expiry_date":   product.ExpiryDate,
			"deleted":       product.Deleted,
			"updated_at":    product.UpdatedAt,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
	return err
}
