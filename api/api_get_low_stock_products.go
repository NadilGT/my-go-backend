package api

import (
	"employee-crud/dao"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetLowStockProductsHandler handles the API request to get top 10 lowest stock products (Fiber version)
func GetLowStockProductsHandler(c *fiber.Ctx) error {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	products, err := dao.DB_GetTopLowStockProducts(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}
