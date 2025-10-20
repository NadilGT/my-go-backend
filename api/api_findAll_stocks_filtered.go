package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// FindAllStocksFilteredApi retrieves stocks filtered by status with cursor-based pagination
// Query params:
//   - cursor: optional, for pagination (pass the next_cursor from previous response)
//   - per_page: optional, default 15, allowed values: 15, 25, 50
//   - status: required, allowed values: "low", "average", "good"
func FindAllStocksFilteredApi(c *fiber.Ctx) error {
	// Get cursor and per_page parameters
	cursor := c.Query("cursor", "")
	perPageStr := c.Query("per_page", "15")
	statusFilter := c.Query("status", "")

	// Validate status filter
	if statusFilter == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "status parameter is required. Allowed values: low, average, good",
		})
	}

	// Normalize status filter
	var minQty, maxQty int
	var statusLabel string

	switch statusFilter {
	case "low":
		minQty = 0
		maxQty = 9 // StockQty < 10
		statusLabel = "Low Stock"
	case "average":
		minQty = 10
		maxQty = 24 // StockQty >= 10 and < 25
		statusLabel = "Average Stock"
	case "good":
		minQty = 25
		maxQty = -1 // StockQty >= 25 (no upper limit)
		statusLabel = "Good Stock"
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status value. Allowed values: low, average, good",
		})
	}

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

	// Use cursor-based pagination with filtering for optimal performance
	stocks, nextCursor, hasMore, err := dao.DB_FindAllStocksFilteredCursorPaginated(perPage, cursor, minQty, maxQty)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get total count for filtered results (optional)
	totalCount, _ := dao.DB_GetStocksCountFiltered(minQty, maxQty)

	response := fiber.Map{
		"data":        stocks,
		"per_page":    perPage,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
		"total_count": totalCount,
		"filter":      statusLabel,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// FindAllStocksFilteredLightweightApi is a lightweight version without total count
// Query params:
//   - cursor: optional, for pagination
//   - per_page: optional, default 15, allowed values: 15, 25, 50
//   - status: required, allowed values: "low", "average", "good"
func FindAllStocksFilteredLightweightApi(c *fiber.Ctx) error {
	// Get cursor and per_page parameters
	cursor := c.Query("cursor", "")
	perPageStr := c.Query("per_page", "15")
	statusFilter := c.Query("status", "")

	// Validate status filter
	if statusFilter == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "status parameter is required. Allowed values: low, average, good",
		})
	}

	// Normalize status filter
	var minQty, maxQty int
	var statusLabel string

	switch statusFilter {
	case "low":
		minQty = 0
		maxQty = 9
		statusLabel = "Low Stock"
	case "average":
		minQty = 10
		maxQty = 24
		statusLabel = "Average Stock"
	case "good":
		minQty = 25
		maxQty = -1
		statusLabel = "Good Stock"
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status value. Allowed values: low, average, good",
		})
	}

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

	// Use cursor-based pagination with filtering
	stocks, nextCursor, hasMore, err := dao.DB_FindAllStocksFilteredCursorPaginated(perPage, cursor, minQty, maxQty)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	type FilteredResponse struct {
		Data       []dto.Stock `json:"data"`
		PerPage    int         `json:"per_page"`
		NextCursor string      `json:"next_cursor"`
		HasMore    bool        `json:"has_more"`
		Filter     string      `json:"filter"`
	}

	response := FilteredResponse{
		Data:       stocks,
		PerPage:    perPage,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Filter:     statusLabel,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
