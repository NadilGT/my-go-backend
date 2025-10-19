package utils

import (
	"employee-crud/dao"
	"log"
	"time"
)

// StartDailyReportScheduler starts a background job that saves daily reports
// It runs every hour and checks if it's time to save yesterday's report
func StartDailyReportScheduler() {
	// Sri Lanka timezone
	sriLankaLoc := time.FixedZone("Asia/Colombo", 5*3600+30*60)

	go func() {
		ticker := time.NewTicker(1 * time.Hour) // Check every hour
		defer ticker.Stop()

		log.Println("Daily Report Scheduler started")

		for {
			select {
			case <-ticker.C:
				now := time.Now().In(sriLankaLoc)

				// Check if it's past midnight (between 00:00 and 01:00)
				// This ensures we save yesterday's report after the day is complete
				if now.Hour() == 0 {
					// Get yesterday's date
					yesterday := now.AddDate(0, 0, -1)
					log.Printf("Attempting to save daily report for: %s\n", yesterday.Format("2006-01-02"))

					// Get the sales summary for yesterday
					summary, err := dao.GetDailySalesSummary(yesterday)
					if err != nil {
						log.Printf("Error getting daily sales summary for %s: %v\n", yesterday.Format("2006-01-02"), err)
						continue
					}

					// Save the report
					err = dao.SaveDailyReport(summary)
					if err != nil {
						log.Printf("Error saving daily report for %s: %v\n", yesterday.Format("2006-01-02"), err)
					} else {
						log.Printf("Successfully saved daily report for %s\n", yesterday.Format("2006-01-02"))
					}
				}

				// Additionally, check for expired reports cleanup at 2 AM
				if now.Hour() == 2 {
					log.Println("Running expired reports cleanup...")
					count, err := dao.DeleteExpiredReports()
					if err != nil {
						log.Printf("Error deleting expired reports: %v\n", err)
					} else if count > 0 {
						log.Printf("Deleted %d expired reports\n", count)
					}
				}
			}
		}
	}()
}

// SaveMissingReports checks and saves reports for any missing dates in the past 7 days
// This is useful for recovering from downtime
func SaveMissingReports() {
	sriLankaLoc := time.FixedZone("Asia/Colombo", 5*3600+30*60)
	today := time.Now().In(sriLankaLoc)

	// Check last 7 days
	for i := 1; i <= 7; i++ {
		checkDate := today.AddDate(0, 0, -i)

		// Check if report already exists
		_, err := dao.GetDailyReportByDate(checkDate)
		if err != nil {
			// Report doesn't exist, try to create it
			log.Printf("Missing report detected for %s, attempting to create...\n", checkDate.Format("2006-01-02"))

			summary, err := dao.GetDailySalesSummary(checkDate)
			if err != nil {
				log.Printf("Error getting sales summary for %s: %v\n", checkDate.Format("2006-01-02"), err)
				continue
			}

			err = dao.SaveDailyReport(summary)
			if err != nil {
				log.Printf("Error saving report for %s: %v\n", checkDate.Format("2006-01-02"), err)
			} else {
				log.Printf("Successfully saved missing report for %s\n", checkDate.Format("2006-01-02"))
			}
		}
	}
}
