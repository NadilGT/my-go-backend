package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func GetProductsBySupplierApi(c *fiber.Ctx) error {
	supplierId := c.Query("supplierId")

	if supplierId == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "supplierId is required", nil)
	}

	products, err := dao.DB_FindProductsBySupplierID(supplierId)
	if err != nil {
		return utils.NewCustomError(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(products)
}
