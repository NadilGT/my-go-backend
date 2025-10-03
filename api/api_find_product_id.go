package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindProductByID(c *fiber.Ctx) error {
	productId := c.Query("productId")

	if productId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	product, err := dao.DB_FindProductByID(productId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}
