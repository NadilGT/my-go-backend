package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

// GetTotalStockQuantityApi retrieves the sum of all stockQty in the Stocks collection
// This calculates the total quantity of all products in stock
func GetTotalStockQuantityApi(c *fiber.Ctx) error {
	// Calculate total stock quantity
	totalQty, err := dao.DB_CalculateTotalStockQuantity()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate total stock quantity",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_stock_quantity": totalQty,
		"message":              "Total stock quantity calculated successfully",
	})
}
