package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

// CleanupOrphanedStocksApi removes stock entries with null or empty batchId
// This is a maintenance endpoint to fix data integrity issues
func CleanupOrphanedStocksApi(c *fiber.Ctx) error {
	deletedCount, err := dao.DB_CleanupOrphanedStocks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Orphaned stock entries cleaned up successfully",
		"deleted_count": deletedCount,
	})
}

// ValidateStockIntegrityApi checks for inconsistencies in stock data
// Returns a report of any issues found
func ValidateStockIntegrityApi(c *fiber.Ctx) error {
	report, err := dao.DB_ValidateStockIntegrity()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(report)
}
