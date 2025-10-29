package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB_FindTopExpiringStocksNext7Days returns top N stocks expiring in next 7 days, sorted by stock quantity desc
func DB_FindTopExpiringStocksNext7Days(topN int, now, sevenDaysLater time.Time) ([]ProductWithStockInfo, error) {
	collection := dbConfigs.DATABASE.Collection("Products")
	ctx := context.Background()

	// Only non-deleted products, with at least one batch expiring in next 7 days
	// (filter variable removed; not used)

	// Find products with batches expiring in next 7 days
	filterBatches := bson.M{
		"deleted": false,
		"batches.expiry_date": bson.M{
			"$gt": now,
			"$lt": sevenDaysLater,
		},
	}
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{
		"productId":       1,
		"name":            1,
		"stockQty":        1,
		"batches":         1,
		"expiry_date":     1,
		"productStockQty": 1,
		"hasBatches":      1,
		"batchCount":      1,
		"productStatus":   1,
		"created_at":      1,
		"updated_at":      1,
	})
	cursor, err := collection.Find(ctx, filterBatches, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var products []struct {
		ProductId  string     `bson:"productId"`
		Name       string     `bson:"name"`
		StockQty   int        `bson:"stockQty"`
		ExpiryDate *time.Time `bson:"expiry_date"`
		Batches    []struct {
			BatchId    string     `bson:"batchId"`
			StockQty   int        `bson:"stockQty"`
			ExpiryDate *time.Time `bson:"expiry_date"`
		} `bson:"batches"`
	}
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	// Find products WITHOUT batches but with expiry_date in next 7 days
	filterNoBatches := bson.M{
		"deleted": false,
		"$or": []bson.M{
			{"batches": bson.M{"$size": 0}},
			{"batches": bson.M{"$exists": false}},
		},
		"expiry_date": bson.M{
			"$gt": now,
			"$lt": sevenDaysLater,
		},
	}
	cursor2, err := collection.Find(ctx, filterNoBatches, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor2.Close(ctx)
	var productsNoBatches []struct {
		ProductId  string     `bson:"productId"`
		Name       string     `bson:"name"`
		StockQty   int        `bson:"stockQty"`
		ExpiryDate *time.Time `bson:"expiry_date"`
	}
	if err := cursor2.All(ctx, &productsNoBatches); err != nil {
		return nil, err
	}

	// Flatten to batch level, filter batches expiring in next 7 days
	var expiringStocks []ProductWithStockInfo
	for _, p := range products {
		if len(p.Batches) == 0 {
			// Product with no batches, but caught by first query (shouldn't happen, but just in case)
			if p.ExpiryDate != nil && p.ExpiryDate.After(now) && p.ExpiryDate.Before(sevenDaysLater) {
				stockInfo := ProductWithStockInfo{
					ProductId:       p.ProductId,
					Name:            p.Name,
					StockQty:        p.StockQty,
					ProductStockQty: p.StockQty,
					ExpiryDate:      p.ExpiryDate,
					BatchId:         "",
					HasBatches:      false,
					BatchCount:      0,
				}
				stockInfo.CalculateProductStatus(p.StockQty)
				expiringStocks = append(expiringStocks, stockInfo)
			}
		} else {
			for _, b := range p.Batches {
				if b.ExpiryDate != nil && b.ExpiryDate.After(now) && b.ExpiryDate.Before(sevenDaysLater) {
					stockInfo := ProductWithStockInfo{
						ProductId:       p.ProductId,
						Name:            p.Name,
						StockQty:        b.StockQty,
						ProductStockQty: p.StockQty,
						ExpiryDate:      b.ExpiryDate,
						BatchId:         b.BatchId,
						HasBatches:      true,
						BatchCount:      len(p.Batches),
					}
					stockInfo.CalculateProductStatus(p.StockQty)
					expiringStocks = append(expiringStocks, stockInfo)
				}
			}
		}
	}
	for _, p := range productsNoBatches {
		if p.ExpiryDate != nil && p.ExpiryDate.After(now) && p.ExpiryDate.Before(sevenDaysLater) {
			stockInfo := ProductWithStockInfo{
				ProductId:       p.ProductId,
				Name:            p.Name,
				StockQty:        p.StockQty,
				ProductStockQty: p.StockQty,
				ExpiryDate:      p.ExpiryDate,
				BatchId:         "",
				HasBatches:      false,
				BatchCount:      0,
			}
			stockInfo.CalculateProductStatus(p.StockQty)
			expiringStocks = append(expiringStocks, stockInfo)
		}
	}

	// Sort by StockQty descending
	if len(expiringStocks) > 1 {
		for i := 0; i < len(expiringStocks)-1; i++ {
			for j := 0; j < len(expiringStocks)-i-1; j++ {
				if expiringStocks[j].StockQty < expiringStocks[j+1].StockQty {
					expiringStocks[j], expiringStocks[j+1] = expiringStocks[j+1], expiringStocks[j]
				}
			}
		}
	}

	// Limit to top N
	if len(expiringStocks) > topN {
		expiringStocks = expiringStocks[:topN]
	}

	return expiringStocks, nil
}
