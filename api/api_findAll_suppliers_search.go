package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func FindAllSuppliersSearch(c *fiber.Ctx) error {
	search := c.Query("search")

	brands, err := dao.DB_FindAllSuppliers("")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	filtered := []dto.Supplier{}
	for _, s := range brands {
		if search != "" {
			if !strings.Contains(strings.ToLower(s.Name), strings.ToLower(search)) {
				continue
			}
		}
		filtered = append(filtered, s)
	}

	return c.Status(fiber.StatusOK).JSON(filtered)
}
