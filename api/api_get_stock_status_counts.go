package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetStockStatusCountsApi retrieves the count of stocks grouped by status
// Returns counts for Low Stock, Average Stock, Good Stock, and Total
// This API scans ALL stocks across ALL pages efficiently using MongoDB aggregation
// OPTIMIZED: Uses 30-second cache to reduce database load for frequent requests
func GetStockStatusCountsApi(c *fiber.Ctx) error {
	cacheKey := "stock_status_counts"

	// Try to get from cache first
	if cached, found := utils.MetricsCache.Get(cacheKey); found {
		if counts, ok := cached.(*dao.StockStatusCounts); ok {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"low_stock":     counts.LowStock,
				"average_stock": counts.AverageStock,
				"good_stock":    counts.GoodStock,
				"total":         counts.Total,
				"message":       "Stock status counts retrieved successfully (cached)",
				"cached":        true,
			})
		}
	}

	// Get stock counts by status from database
	counts, err := dao.DB_GetStockStatusCounts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve stock status counts",
		})
	}

	// Cache for 30 seconds
	utils.MetricsCache.Set(cacheKey, counts, 30*time.Second)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"low_stock":     counts.LowStock,
		"average_stock": counts.AverageStock,
		"good_stock":    counts.GoodStock,
		"total":         counts.Total,
		"message":       "Stock status counts retrieved successfully",
		"cached":        false,
	})
}
