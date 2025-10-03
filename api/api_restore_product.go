package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func RestoreProductApi(c *fiber.Ctx) error {
	productId := c.Query("productId")
	categoryId := c.Query("categoryId")
	brandId := c.Query("brandId")
	subCategoryId := c.Query("subCategoryId")

	if productId == "" || categoryId == "" || brandId == "" || subCategoryId == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "ProductId, CategoryId, BrandId and SubCategoryId are required", nil)
	}

	if err := dao.DB_RestoreProductByID(productId, categoryId, brandId, subCategoryId); err != nil {
		return utils.NewCustomError(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Product restored successfully",
		"productId":     productId,
		"categoryId":    categoryId,
		"brandId":       brandId,
		"subCategoryId": subCategoryId,
	})
}
