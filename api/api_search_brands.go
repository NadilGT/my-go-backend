package api

import (
	"employee-crud/dao"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// SearchBrandsApi searches for brands by name with optional limit
func SearchBrandsApi(c *fiber.Ctx) error {
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

	brands, err := dao.DB_SearchBrands(searchTerm, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": brands,
	})
}
