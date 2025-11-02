package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindAllSuppliers(c *fiber.Ctx) error {
	status := c.Query("status")
	suppliers, err := dao.DB_FindAllSuppliers(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(suppliers)
}
