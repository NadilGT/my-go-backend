package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetProductsByBrandApi(c *fiber.Ctx) error {
	brandId := c.Query("brandId")

	if brandId == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "brandId is required", nil)
	}

	// Get cursor and per_page parameters
	cursor := c.Query("cursor", "")
	perPageStr := c.Query("per_page", "15")

	// Parse per_page
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		perPage = 15
	}

	// Validate per_page values (only allow 15, 25, 50)
	switch perPage {
	case 15, 25, 50:
		// Valid per_page value
	default:
		perPage = 15 // Default to 15 if invalid value provided
	}

	// Use cursor-based pagination for optimal performance
	products, nextCursor, hasMore, err := dao.DB_FindProductsByBrandCursorPaginated(brandId, perPage, cursor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := utils.CursorPaginatedResponse{
		Data:       products,
		PerPage:    perPage,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
