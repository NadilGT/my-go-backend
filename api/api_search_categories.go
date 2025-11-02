package api

import (
	"employee-crud/dao"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// SearchCategoriesApi searches for categories by name with optional limit
func SearchCategoriesApi(c *fiber.Ctx) error {
	searchTerm := c.Query("q", "")
	limitStr := c.Query("limit", "20")

	// Convert limit to integer
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	// Cap the limit at 100 to prevent excessive results
	if limit > 100 {
		limit = 100
	}

	categories, err := dao.DB_SearchCategories(searchTerm, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": categories,
	})
}
