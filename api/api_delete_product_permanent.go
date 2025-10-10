package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteProductPermanentApi(c *fiber.Ctx) error {
	productId := c.Query("productId")

	if productId == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "Product ID is required", nil)
	}

	if err := dao.DB_PermanentDeleteProductByID(productId); err != nil {
		return utils.NewCustomError(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product permanently deleted successfully",
	})
}
