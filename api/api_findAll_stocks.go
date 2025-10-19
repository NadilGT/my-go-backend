package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// FindAllStocksApi retrieves all stocks with cursor-based pagination
// Query params:
//   - cursor: optional, for pagination (pass the next_cursor from previous response)
//   - per_page: optional, default 15, allowed values: 15, 25, 50
func FindAllStocksApi(c *fiber.Ctx) error {
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

	// Use cursor-based pagination for optimal performance with large datasets
	stocks, nextCursor, hasMore, err := dao.DB_FindAllStocksCursorPaginated(perPage, cursor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get total count (optional, can be removed for better performance)
	totalCount, _ := dao.DB_GetStocksCount()

	response := fiber.Map{
		"data":        stocks,
		"per_page":    perPage,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
		"total_count": totalCount,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// FindAllStocksLightweightApi is a lightweight version without total count
// Use this for better performance when you don't need the total count
func FindAllStocksLightweightApi(c *fiber.Ctx) error {
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
	stocks, nextCursor, hasMore, err := dao.DB_FindAllStocksCursorPaginated(perPage, cursor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := utils.CursorPaginatedResponse{
		Data:       stocks,
		PerPage:    perPage,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
