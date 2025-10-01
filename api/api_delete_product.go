package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteProductApi(c *fiber.Ctx) error {
	id := c.Query("productId")

	if id == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "Product ID is required", nil)
	}

	if err := dao.DB_DeleteProductByID(id); err != nil {
		return utils.NewCustomError(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
