package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetTotalStockQuantityApi retrieves the sum of all stockQty in the Stocks collection
// This calculates the total quantity of all products in stock
// OPTIMIZED: Uses 30-second cache to reduce database load for frequent requests
func GetTotalStockQuantityApi(c *fiber.Ctx) error {
	cacheKey := "total_stock_quantity"

	// Try to get from cache first
	if cached, found := utils.MetricsCache.Get(cacheKey); found {
		if totalQty, ok := cached.(int64); ok {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"total_stock_quantity": totalQty,
				"message":              "Total stock quantity calculated successfully (cached)",
				"cached":               true,
			})
		}
	}

	// Calculate total stock quantity from database
	totalQty, err := dao.DB_CalculateTotalStockQuantity()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate total stock quantity",
		})
	}

	// Cache for 30 seconds
	utils.MetricsCache.Set(cacheKey, totalQty, 30*time.Second)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_stock_quantity": totalQty,
		"message":              "Total stock quantity calculated successfully",
		"cached":               false,
	})
}
