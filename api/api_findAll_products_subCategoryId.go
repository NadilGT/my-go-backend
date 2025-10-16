package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllProductsBySubCategoryApi(c *fiber.Ctx) error {
	categoryId := c.Query("subCategoryId")

	if categoryId == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "SubCategoryId is required", nil)
	}

	products, err := dao.DB_FindProductsBySubCategory(categoryId)
	if err != nil {
		return utils.NewCustomError(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if products == nil {
		products = []dto.Product{}
	}

	return c.Status(fiber.StatusOK).JSON(products)
}
