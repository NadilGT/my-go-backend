package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

// GetSalesReportsApi returns all sales reports (available for 30 days)
func GetSalesReportsApi(c *fiber.Ctx) error {
	reports, err := dao.DB_FindAllSalesReports()
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(reports)
}

// GetSalesReportByDateApi returns a specific sales report by date
func GetSalesReportByDateApi(c *fiber.Ctx) error {
	date := c.Query("date") // Expected format: 2024-01-15
	if date == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "date parameter is required (format: YYYY-MM-DD)")
	}

	report, err := dao.DB_FindSalesReportByDate(date)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Report not found for the specified date")
	}

	return c.JSON(report)
}
