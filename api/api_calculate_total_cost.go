package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func CalculateTotalCostPrice(c *fiber.Ctx) error {
	total, err := dao.DB_CalculateTotalCostPrice()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_cost_value": total,
	})
}
