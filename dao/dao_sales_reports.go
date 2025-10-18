package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB_GenerateDailySalesReport generates a daily sales report from existing sales data
// This should be called by a scheduled job before deleting old sales
func DB_GenerateDailySalesReport(reportDate time.Time) error {
	collection := dbConfigs.DATABASE.Collection("Sales")
	ctx := context.Background()

	// Normalize the date to start of day
	startOfDay := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, reportDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Find all sales for the specified date
	cursor, err := collection.Find(
		ctx,
		bson.M{
			"createdAt": bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			},
			"deleted": false,
		},
	)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var sales []dto.Sale
	if err := cursor.All(ctx, &sales); err != nil {
		return err
	}

	// If no sales for the day, still create a report with zeros
	if len(sales) == 0 {
		return createEmptyReport(reportDate)
	}

	// Calculate aggregate data
	report := calculateDailyReport(sales, reportDate)

	// Save the report
	return DB_CreateSalesReport(ctx, &report)
}

func calculateDailyReport(sales []dto.Sale, reportDate time.Time) dto.SaleReport {
	report := dto.SaleReport{
		ReportDate:  reportDate,
		GeneratedAt: time.Now().UTC(),
		GeneratedBy: "System",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Product tracking for top products
	productMap := make(map[string]*dto.TopProductReport)

	for _, sale := range sales {
		report.TotalSales++
		report.TotalRevenue += sale.GrandTotal
		report.TotalDiscount += sale.TotalDiscount
		report.TotalTax += sale.TaxAmount

		// Payment method breakdown
		switch sale.PaymentMethod {
		case "cash":
			report.CashSales++
			report.CashRevenue += sale.GrandTotal
		case "card":
			report.CardSales++
			report.CardRevenue += sale.GrandTotal
		case "transfer":
			report.TransferSales++
			report.TransferRevenue += sale.GrandTotal
		}

		// Customer info tracking
		if sale.CustomerName != "" {
			report.SalesWithCustomer++
		} else {
			report.SalesWithoutCustomer++
		}

		// Track products
		for _, item := range sale.Items {
			report.TotalItemsSold += item.Quantity

			if existing, exists := productMap[item.ProductId]; exists {
				existing.QuantitySold += item.Quantity
				existing.Revenue += item.TotalPrice
			} else {
				productMap[item.ProductId] = &dto.TopProductReport{
					ProductId:    item.ProductId,
					ProductName:  item.ProductName,
					QuantitySold: item.Quantity,
					Revenue:      item.TotalPrice,
				}
			}
		}
	}

	// Convert product map to slice and get top products
	topProducts := []dto.TopProductReport{}
	for _, product := range productMap {
		topProducts = append(topProducts, *product)
	}

	// Sort by quantity sold (simple bubble sort for small lists)
	for i := 0; i < len(topProducts); i++ {
		for j := i + 1; j < len(topProducts); j++ {
			if topProducts[j].QuantitySold > topProducts[i].QuantitySold {
				topProducts[i], topProducts[j] = topProducts[j], topProducts[i]
			}
		}
	}

	// Keep only top 10
	if len(topProducts) > 10 {
		topProducts = topProducts[:10]
	}

	report.TopProducts = topProducts

	return report
}

func createEmptyReport(reportDate time.Time) error {
	report := dto.SaleReport{
		ReportDate:           reportDate,
		TotalSales:           0,
		TotalRevenue:         0,
		TotalDiscount:        0,
		TotalTax:             0,
		TotalItemsSold:       0,
		CashSales:            0,
		CardSales:            0,
		TransferSales:        0,
		CashRevenue:          0,
		CardRevenue:          0,
		TransferRevenue:      0,
		SalesWithCustomer:    0,
		SalesWithoutCustomer: 0,
		GeneratedAt:          time.Now().UTC(),
		GeneratedBy:          "System",
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	return DB_CreateSalesReport(context.Background(), &report)
}

func DB_CreateSalesReport(ctx context.Context, report *dto.SaleReport) error {
	collection := dbConfigs.DATABASE.Collection("SalesReports")

	// Generate report ID
	id, err := GenerateId(ctx, "SalesReports", "RPT")
	if err != nil {
		return err
	}
	report.ReportId = id

	_, err = collection.InsertOne(ctx, report)
	return err
}

func DB_FindAllSalesReports() ([]dto.SaleReport, error) {
	collection := dbConfigs.DATABASE.Collection("SalesReports")
	ctx := context.Background()

	// Sort by report date descending (newest first)
	opts := options.Find().SetSort(bson.D{{Key: "reportDate", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []dto.SaleReport
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

func DB_FindSalesReportByDate(dateStr string) (*dto.SaleReport, error) {
	collection := dbConfigs.DATABASE.Collection("SalesReports")
	ctx := context.Background()

	// Parse date string
	reportDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	// Normalize to start of day
	startOfDay := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, reportDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var report dto.SaleReport
	err = collection.FindOne(
		ctx,
		bson.M{
			"reportDate": bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			},
		},
	).Decode(&report)

	if err != nil {
		return nil, err
	}

	return &report, nil
}

// DB_DeleteReportsOlderThan30Days deletes sales reports that are older than 30 days
func DB_DeleteReportsOlderThan30Days() (int64, error) {
	collection := dbConfigs.DATABASE.Collection("SalesReports")
	ctx := context.Background()

	thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)

	result, err := collection.DeleteMany(
		ctx,
		bson.M{
			"reportDate": bson.M{
				"$lt": thirtyDaysAgo,
			},
		},
	)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
