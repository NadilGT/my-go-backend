package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func FindAllProductsSearch(c *fiber.Ctx) error {
	search := c.Query("search")

	products, err := dao.DB_FindAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	filtered := []dto.Product{}
	for _, p := range products {
		if search != "" {
			if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) {
				continue
			}
		}
		filtered = append(filtered, p)
	}

	return c.Status(fiber.StatusOK).JSON(filtered)
}
