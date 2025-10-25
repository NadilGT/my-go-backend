package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// FindAllProductsWithStockApi retrieves all products with their stock information
// This includes both products with batches and products without batches
// Query params:
//   - cursor: optional, for pagination (pass the next_cursor from previous response)
//   - per_page: optional, default 15, allowed values: 15, 25, 50
func FindAllProductsWithStockApi(c *fiber.Ctx) error {
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

	// Get products with stock
	productsWithStock, nextCursor, hasMore, err := dao.DB_FindAllProductsWithStock(perPage, cursor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get total count (optional)
	totalCount, _ := dao.DB_GetProductsWithStockCount()

	response := fiber.Map{
		"data":        productsWithStock,
		"per_page":    perPage,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
		"total_count": totalCount,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// FindAllProductsWithStockLiteApi is a lightweight version without total count
// Query params:
//   - cursor: optional, for pagination
//   - per_page: optional, default 15, allowed values: 15, 25, 50
func FindAllProductsWithStockLiteApi(c *fiber.Ctx) error {
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

	// Get products with stock
	productsWithStock, nextCursor, hasMore, err := dao.DB_FindAllProductsWithStock(perPage, cursor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := utils.CursorPaginatedResponse{
		Data:       productsWithStock,
		PerPage:    perPage,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
