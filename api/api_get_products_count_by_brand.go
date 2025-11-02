package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

// GetProductsCountByBrandApi returns the count of products for a specific brand (non-deleted only)
func GetProductsCountByBrandApi(c *fiber.Ctx) error {
	brandId := c.Params("brandId")

	if brandId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "brandId parameter is required",
		})
	}

	count, err := dao.DB_GetProductsCountByBrand(brandId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": count,
	})
}
