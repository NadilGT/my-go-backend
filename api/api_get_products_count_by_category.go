package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

// GetProductsCountByCategoryApi returns the count of products for a specific category (non-deleted only)
func GetProductsCountByCategoryApi(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")

	if categoryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "categoryId parameter is required",
		})
	}

	count, err := dao.DB_GetProductsCountByCategory(categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": count,
	})
}
