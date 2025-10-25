package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

// SyncStocksApi syncs all product stocks to the Stocks collection
// This should be called after initial setup or periodically to ensure stocks are up to date
// For large datasets, this may take some time
func SyncStocksApi(c *fiber.Ctx) error {
	err := dao.DB_SyncStocksFromProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"operation": "Failed",
			"error":     err.Error(),
			"message":   "Failed to sync stocks from products",
		})
	}

	// Invalidate metrics cache after stock sync
	utils.MetricsCache.Delete("stock_status_counts")
	utils.MetricsCache.Delete("total_stock_quantity")

	// Get the total count of synced stocks
	count, _ := dao.DB_GetStocksCount()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"operation":    "Success",
		"message":      "Stocks synchronized successfully",
		"total_stocks": count,
	})
}
