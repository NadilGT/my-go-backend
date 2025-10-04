package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("categoryId %s not found", product.ProductId)
	}

	return nil
}
