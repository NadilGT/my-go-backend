package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func FindAllBrandsSearch(c *fiber.Ctx) error {
	search := c.Query("search")

	brands, err := dao.DB_FindAllBrands()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	filtered := []dto.Brand{}
	for _, p := range brands {
		if search != "" {
			if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) {
				continue
			}
		}
		filtered = append(filtered, p)
	}

	return c.Status(fiber.StatusOK).JSON(filtered)
}
