package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteSupplierApi(c *fiber.Ctx) error {
	id := c.Query("supplierId")

	if id == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "Supplier ID is required", nil)
	}

	if err := dao.DB_DeleteSupplierByID(id); err != nil {
		return utils.NewCustomError(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Supplier deleted successfully",
	})
}
