package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

// GetTotalProducts returns the total number of products in the database
func GetTotalProducts(c *fiber.Ctx) error {
	total, err := dao.DB_GetTotalProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_products": total,
	})
}
