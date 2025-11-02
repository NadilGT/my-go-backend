package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

// GetCategorizedProductsCountApi returns the total count of products that have a categoryId assigned and are not deleted
func GetCategorizedProductsCountApi(c *fiber.Ctx) error {
	count, err := dao.DB_GetCategorizedProductsCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": count,
	})
}
