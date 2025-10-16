package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindAllGRNs(c *fiber.Ctx) error {
	grns, err := dao.DB_FindAllGRNs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(grns)
}
