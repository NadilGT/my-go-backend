package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func DB_UpdateSupplier(ctx context.Context, supplier *dto.Supplier) error {
	collection := dbConfigs.DATABASE.Collection("Suppliers")

	filter := bson.M{"supplierId": supplier.SupplierId}

	update := bson.M{
		"$set": bson.M{
			"name":       supplier.Name,
			"contact":    supplier.Contact,
			"email":      supplier.Email,
			"address":    supplier.Address,
			"updated_at": supplier.UpdatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("supplierId %s not found", supplier.SupplierId)
	}

	return nil
}
