package api

import (
	"employee-crud/dao"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetSavedDailyReportApi retrieves a saved daily report by date
func GetSavedDailyReportApi(c *fiber.Ctx) error {
	// Get date parameter from query (format: YYYY-MM-DD)
	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Date parameter is required. Use format YYYY-MM-DD (e.g., 2025-10-19)",
		})
	}

	// Parse the provided date
	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Use YYYY-MM-DD (e.g., 2025-10-19)",
		})
	}

	// Get the saved report
	report, err := dao.GetDailyReportByDate(targetDate)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No saved report found for the specified date",
		})
	}

	return c.JSON(report)
}

// GetMonthlyReportsApi retrieves all saved daily reports for a specific month
func GetMonthlyReportsApi(c *fiber.Ctx) error {
	// Get year and month parameters
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Year and month parameters are required (e.g., ?year=2025&month=10)",
		})
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid year parameter",
		})
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid month parameter (must be 1-12)",
		})
	}

	// Get the reports
	reports, err := dao.GetDailyReportsByMonth(year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve monthly reports: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"year":         year,
		"month":        month,
		"reportCount":  len(reports),
		"dailyReports": reports,
	})
}
