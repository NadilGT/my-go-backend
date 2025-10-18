package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindSaleByIdApi(c *fiber.Ctx) error {
	saleId := c.Query("saleId")

	if saleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Sale ID is required",
		})
	}

	sale, err := dao.FindSaleBySaleId(saleId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sale not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(sale)
}
