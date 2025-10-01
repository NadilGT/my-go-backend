package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func AssignProductToSupplierApi(c *fiber.Ctx) error {
	supplierId := c.Query("supplierId")
	productId := c.Query("productId")

	if supplierId == "" || productId == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "supplierId and productId are required", nil)
	}

	if err := dao.DB_AssignProductToSupplier(supplierId, productId); err != nil {
		return utils.NewCustomError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product assigned to supplier successfully",
	})
}
