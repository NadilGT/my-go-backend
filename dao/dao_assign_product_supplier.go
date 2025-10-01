package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_AssignProductToSupplier(supplierId, productId string) error {
	supplierColl := dbConfigs.DATABASE.Collection("Suppliers")
	productColl := dbConfigs.DATABASE.Collection("Products")
	assignColl := dbConfigs.DATABASE.Collection("SupplierProducts")
	ctx := context.Background()

	var supplier dto.Supplier
	err := supplierColl.FindOne(ctx, bson.M{"supplierId": supplierId, "deleted": false}).Decode(&supplier)
	if err != nil {
		return errors.New("supplier not found")
	}

	var product dto.Product
	err = productColl.FindOne(ctx, bson.M{"productId": productId}).Decode(&product)
	if err != nil {
		return errors.New("product not found")
	}

	filter := bson.M{
		"supplierId": supplierId,
		"productId":  productId,
	}
	count, err := assignColl.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("product already assigned to this supplier")
	}

	now := time.Now()
	newAssignment := dto.SupplierProduct{
		SupplierID:   supplier.SupplierId,
		SupplierName: supplier.Name,
		ProductID:    product.ProductId,
		ProductName:  product.Name,
		AssignedAt:   now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = assignColl.InsertOne(ctx, newAssignment)
	if err != nil {
		return err
	}

	return nil
}
