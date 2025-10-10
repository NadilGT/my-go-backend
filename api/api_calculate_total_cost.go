package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func CalculateTotalAndExpectedCost(c *fiber.Ctx) error {
	totalCost, expectedCost, err := dao.DB_CalculateTotalAndExpectedCost()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"target_profit": expectedCost - totalCost,
		"sales_target":  expectedCost,
		"total_spend":   totalCost,
	})
}
