package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

// GetStockStatusCountsApi retrieves the count of stocks grouped by status
// Returns counts for Low Stock, Average Stock, Good Stock, and Total
// This API scans ALL stocks across ALL pages efficiently using MongoDB aggregation
func GetStockStatusCountsApi(c *fiber.Ctx) error {
	// Get stock counts by status
	counts, err := dao.DB_GetStockStatusCounts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve stock status counts",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"low_stock":     counts.LowStock,
		"average_stock": counts.AverageStock,
		"good_stock":    counts.GoodStock,
		"total":         counts.Total,
		"message":       "Stock status counts retrieved successfully",
	})
}
