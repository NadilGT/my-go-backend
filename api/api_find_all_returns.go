package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindAllReturnsApi(c *fiber.Ctx) error {
	results, err := dao.GetAllReturns(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch returns"})
	}
	return c.JSON(results)
}
