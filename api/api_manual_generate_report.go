package api

import (
	"employee-crud/functions"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

// ManuallyGenerateReportApi allows manual report generation for a specific date
// Useful for testing or regenerating reports
func ManuallyGenerateReportApi(c *fiber.Ctx) error {
	date := c.Query("date") // Expected format: 2024-01-15
	if date == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "date parameter is required (format: YYYY-MM-DD)")
	}

	if err := functions.ManuallyGenerateReportForDate(date); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Report generated successfully for " + date,
	})
}
