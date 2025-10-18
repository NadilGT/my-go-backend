package functions

import (
	"employee-crud/dao"
	"log"
	"time"
)

// StartSalesCleanupScheduler starts the background job for cleaning up old sales data
// and generating daily reports
func StartSalesCleanupScheduler() {
	// Run immediately on startup to clean any old data
	go runScheduledTasks()

	// Schedule to run every hour
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			runScheduledTasks()
		}
	}()

	log.Println("Sales cleanup scheduler started - runs every hour")
}

func runScheduledTasks() {
	now := time.Now()
	log.Printf("[Scheduler] Running scheduled tasks at %s", now.Format("2006-01-02 15:04:05"))

	// Task 1: Generate report for yesterday (if not already generated)
	yesterday := now.Add(-24 * time.Hour)
	if err := dao.DB_GenerateDailySalesReport(yesterday); err != nil {
		log.Printf("[Scheduler] Error generating report for %s: %v", yesterday.Format("2006-01-02"), err)
	} else {
		log.Printf("[Scheduler] Successfully generated/updated report for %s", yesterday.Format("2006-01-02"))
	}

	// Task 2: Delete sales older than 24 hours
	deletedSales, err := dao.DB_DeleteSalesOlderThan24Hours()
	if err != nil {
		log.Printf("[Scheduler] Error deleting old sales: %v", err)
	} else if deletedSales > 0 {
		log.Printf("[Scheduler] Deleted %d sales older than 24 hours", deletedSales)
	}

	// Task 3: Delete reports older than 30 days
	deletedReports, err := dao.DB_DeleteReportsOlderThan30Days()
	if err != nil {
		log.Printf("[Scheduler] Error deleting old reports: %v", err)
	} else if deletedReports > 0 {
		log.Printf("[Scheduler] Deleted %d reports older than 30 days", deletedReports)
	}

	log.Println("[Scheduler] Scheduled tasks completed")
}

// ManuallyGenerateReportForDate allows manual report generation for a specific date
// Useful for testing or regenerating reports
func ManuallyGenerateReportForDate(dateStr string) error {
	reportDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	return dao.DB_GenerateDailySalesReport(reportDate)
}
