package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func GetProductsByCategoryApi(c *fiber.Ctx) error {
	categoryId := c.Query("categoryId")

	if categoryId == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "CategoryId is required", nil)
	}

	products, err := dao.DB_FindProductsByCategory(categoryId)
	if err != nil {
		return utils.NewCustomError(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(products)
}
