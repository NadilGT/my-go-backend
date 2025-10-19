package dao

import (
	"context"
	"employee-crud/dbConfigs"
	"employee-crud/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveDailyReport saves the daily sales report to the database
func SaveDailyReport(summary *dto.DailySalesSummary) error {
	collection := dbConfigs.DATABASE.Collection("DailyReports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate expiration date (end of the month + 1 day)
	// For example, October report expires on November 1st 00:00:00
	reportDate := summary.ReportDate
	year := reportDate.Year()
	month := reportDate.Month()

	// Get the first day of next month
	var expiresAt time.Time
	if month == 12 {
		expiresAt = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Colombo", 5*3600+30*60))
	} else {
		expiresAt = time.Date(year, month+1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Colombo", 5*3600+30*60))
	}

	// Create document
	report := dto.DailyReportDocument{
		ReportDate:      summary.ReportDate,
		Month:           int(summary.ReportDate.Month()),
		Year:            summary.ReportDate.Year(),
		TotalSales:      summary.TotalSales,
		TotalRevenue:    summary.TotalRevenue,
		TotalDiscount:   summary.TotalDiscount,
		TotalTax:        summary.TotalTax,
		CashSales:       summary.CashSales,
		CardSales:       summary.CardSales,
		CashRevenue:     summary.CashRevenue,
		CardRevenue:     summary.CardRevenue,
		ProductsSold:    summary.ProductsSold,
		TopSellingItems: summary.TopSellingItems,
		CreatedAt:       time.Now().In(time.FixedZone("Asia/Colombo", 5*3600+30*60)),
		ExpiresAt:       expiresAt,
	}

	// Check if report already exists for this date
	filter := bson.M{
		"reportDate": bson.M{
			"$gte": time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, time.FixedZone("Asia/Colombo", 5*3600+30*60)),
			"$lt":  time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day()+1, 0, 0, 0, 0, time.FixedZone("Asia/Colombo", 5*3600+30*60)),
		},
	}

	// Use upsert to either insert or update
	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": report}

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetDailyReportByDate retrieves a daily report by date
func GetDailyReportByDate(date time.Time) (*dto.DailyReportDocument, error) {
	collection := dbConfigs.DATABASE.Collection("DailyReports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"reportDate": bson.M{
			"$gte": time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.FixedZone("Asia/Colombo", 5*3600+30*60)),
			"$lt":  time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, time.FixedZone("Asia/Colombo", 5*3600+30*60)),
		},
	}

	var report dto.DailyReportDocument
	err := collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

// GetDailyReportsByMonth retrieves all daily reports for a specific month
func GetDailyReportsByMonth(year int, month int) ([]dto.DailyReportDocument, error) {
	collection := dbConfigs.DATABASE.Collection("DailyReports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"year":  year,
		"month": month,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []dto.DailyReportDocument
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

// DeleteExpiredReports manually deletes reports that have passed their expiration date
// This is a backup function in case TTL index doesn't work properly
func DeleteExpiredReports() (int64, error) {
	collection := dbConfigs.DATABASE.Collection("DailyReports")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now().In(time.FixedZone("Asia/Colombo", 5*3600+30*60))
	filter := bson.M{
		"expiresAt": bson.M{"$lte": now},
	}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
