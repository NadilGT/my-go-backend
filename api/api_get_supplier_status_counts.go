package api

import (
	"context"
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func GetSupplierStatusCounts(c *fiber.Ctx) error {
	active, inactive, err := dao.DB_GetSupplierStatusCounts(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"active":   active,
		"inactive": inactive,
	})
}
