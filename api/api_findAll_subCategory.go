package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindAllSubCategory(c *fiber.Ctx) error {
	brands, err := dao.DB_FindAllSubCategory()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(brands)
}
