package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetDailySalesSummary retrieves sales summary for a specific date
func GetDailySalesSummary(targetDate time.Time) (*dto.DailySalesSummary, error) {
	collection := dbConfigs.DATABASE.Collection("Sales")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Set the date range for the target date (start of day to end of day in UTC)
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Query filter for the specific date
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	// Find all sales for the date
	cursor, err := collection.Find(ctx, filter, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []dto.Sale
	if err := cursor.All(ctx, &sales); err != nil {
		return nil, err
	}

	// Calculate summary
	summary := &dto.DailySalesSummary{
		ReportDate:   targetDate,
		ProductsSold: make([]dto.ProductSoldSummary, 0),
	}

	// Map to aggregate product sales
	productMap := make(map[string]*dto.ProductSoldSummary)

	// Process each sale
	for _, sale := range sales {
		summary.TotalSales++
		summary.TotalRevenue += sale.Total
		summary.TotalDiscount += sale.Discount
		summary.TotalTax += sale.Tax

		// Count payment methods
		if sale.PaymentMethod == "cash" {
			summary.CashSales++
			summary.CashRevenue += sale.Total
		} else if sale.PaymentMethod == "card" {
			summary.CardSales++
			summary.CardRevenue += sale.Total
		}

		// Aggregate product sales
		for _, item := range sale.Items {
			if existing, exists := productMap[item.ProductID]; exists {
				existing.Quantity += item.Quantity
				existing.TotalAmount += item.TotalPrice
			} else {
				productMap[item.ProductID] = &dto.ProductSoldSummary{
					ProductID:   item.ProductID,
					ProductName: item.ProductName,
					Quantity:    item.Quantity,
					UnitPrice:   item.UnitPrice,
					TotalAmount: item.TotalPrice,
				}
			}
		}
	}

	// Convert map to slice
	for _, product := range productMap {
		summary.ProductsSold = append(summary.ProductsSold, *product)
	}

	// Sort products by total amount (descending)
	sort.Slice(summary.ProductsSold, func(i, j int) bool {
		return summary.ProductsSold[i].TotalAmount > summary.ProductsSold[j].TotalAmount
	})

	// Get top 10 selling items
	topCount := 10
	if len(summary.ProductsSold) < topCount {
		topCount = len(summary.ProductsSold)
	}
	summary.TopSellingItems = summary.ProductsSold[:topCount]

	return summary, nil
}
