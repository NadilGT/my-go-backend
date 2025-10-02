package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func GetProductsByBrandApi(c *fiber.Ctx) error {
	brandId := c.Query("brandId")

	if brandId == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "brandId is required", nil)
	}

	products, err := dao.DB_FindProductsByBrand(brandId)
	if err != nil {
		return utils.NewCustomError(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if products == nil {
		products = []dto.Product{}
	}

	return c.Status(fiber.StatusOK).JSON(products)
}
