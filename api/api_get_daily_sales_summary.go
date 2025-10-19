package api

import (
	"employee-crud/dao"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetDailySalesSummaryApi(c *fiber.Ctx) error {
	// Get date parameter from query (format: YYYY-MM-DD)
	// If not provided, use today's date
	dateStr := c.Query("date")
	var targetDate time.Time
	var err error

	if dateStr == "" {
		// Use today's date in Sri Lanka timezone (UTC+5:30)
		sriLankaLoc := time.FixedZone("Asia/Colombo", 5*3600+30*60)
		targetDate = time.Now().In(sriLankaLoc)
	} else {
		// Parse the provided date
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Use YYYY-MM-DD (e.g., 2025-10-19)",
			})
		}
	}

	// Get sales summary for the date
	summary, err := dao.GetDailySalesSummary(targetDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve sales summary: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Sales summary retrieved successfully",
		"data":    summary,
	})
}
