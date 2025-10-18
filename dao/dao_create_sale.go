package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DB_CreateSale(ctx context.Context, sale *dto.Sale) error {
	collection := dbConfigs.DATABASE.Collection("Sales")

	_, err := collection.InsertOne(ctx, sale)
	if err != nil {
		return err
	}

	return nil
}

func DB_UpdateProductStockAfterSale(ctx context.Context, items []dto.SaleItem) error {
	collection := dbConfigs.DATABASE.Collection("Products")

	for _, item := range items {
		// Get current product to check stock
		var product dto.Product
		err := collection.FindOne(ctx, bson.M{"productId": item.ProductId}).Decode(&product)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return errors.New("product not found: " + item.ProductId)
			}
			return err
		}

		// Check if enough stock is available
		if product.StockQty < item.Quantity {
			return errors.New("insufficient stock for product: " + item.ProductName)
		}

		// Update stock quantity
		newStockQty := product.StockQty - item.Quantity
		_, err = collection.UpdateOne(
			ctx,
			bson.M{"productId": item.ProductId},
			bson.M{"$set": bson.M{"stockQty": newStockQty}},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func DB_FindSaleById(saleId string) (*dto.Sale, error) {
	collection := dbConfigs.DATABASE.Collection("Sales")
	ctx := context.Background()

	var sale dto.Sale
	err := collection.FindOne(ctx, bson.M{"saleId": saleId, "deleted": false}).Decode(&sale)
	if err != nil {
		return nil, err
	}

	return &sale, nil
}

func DB_FindAllSales() ([]dto.Sale, error) {
	collection := dbConfigs.DATABASE.Collection("Sales")
	ctx := context.Background()

	cursor, err := collection.Find(ctx, bson.M{"deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []dto.Sale
	if err := cursor.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, nil
}

// DB_FindSalesLast24Hours returns all sales from the last 24 hours
func DB_FindSalesLast24Hours() ([]interface{}, error) {
	collection := dbConfigs.DATABASE.Collection("Sales")
	ctx := context.Background()

	// Calculate 24 hours ago
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

	// Find sales created in the last 24 hours
	cursor, err := collection.Find(
		ctx,
		bson.M{
			"deleted": false,
			"createdAt": bson.M{
				"$gte": twentyFourHoursAgo,
			},
		},
		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []interface{}
	if err := cursor.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, nil
}

// DB_DeleteSalesOlderThan24Hours deletes sales that are older than 24 hours
// This should be called by a scheduled job/cron
func DB_DeleteSalesOlderThan24Hours() (int64, error) {
	collection := dbConfigs.DATABASE.Collection("Sales")
	ctx := context.Background()

	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

	result, err := collection.DeleteMany(
		ctx,
		bson.M{
			"createdAt": bson.M{
				"$lt": twentyFourHoursAgo,
			},
		},
	)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
