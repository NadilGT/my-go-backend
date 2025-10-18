package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindAllSalesApi(c *fiber.Ctx) error {
	sales, err := dao.FindAllSales(0, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sales",
		})
	}

	return c.Status(fiber.StatusOK).JSON(sales)
}
