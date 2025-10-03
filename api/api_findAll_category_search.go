package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func FindAllCategoriesSearchApi(c *fiber.Ctx) error {
	search := c.Query("search")

	categories, err := dao.DB_FindAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	filtered := []dto.Category{}
	for _, p := range categories {
		if search != "" {
			if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) {
				continue
			}
		}
		filtered = append(filtered, p)
	}

	return c.Status(fiber.StatusOK).JSON(filtered)
}
