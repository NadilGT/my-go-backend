package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func GetTotalGRNsCount(c *fiber.Ctx) error {
	count, err := dao.DB_CountTotalGRNs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_grns": count,
	})
}
