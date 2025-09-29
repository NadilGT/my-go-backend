package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindAllProducts(c *fiber.Ctx) error {
	brands, err := dao.DB_FindAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(brands)
}
