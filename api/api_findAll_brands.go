package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FindAllBrands(c *fiber.Ctx) error {
	// Get pagination parameters
	pageStr := c.Query("page", "1")
	perPageStr := c.Query("per_page", "15")

	// Convert page to integer
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Convert per_page to integer and validate allowed values
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

	// Always use paginated version for consistent response format
	brands, total, err := dao.DB_FindAllBrandsPaginated(page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get product counts for all brands in this page
	brandIds := make([]string, len(brands))
	for i, brand := range brands {
		brandIds[i] = brand.BrandId
	}

	if len(brandIds) > 0 {
		productCounts, err := dao.DB_GetProductCountsForBrands(brandIds)
		if err != nil {
			// Log error but continue without product counts
			productCounts = make(map[string]int64)
		}

		// Add product counts to brands
		for i := range brands {
			brands[i].ProductCount = productCounts[brands[i].BrandId]
		}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	// Validate if requested page exists
	if page > totalPages && total > 0 {
		// Return empty data for pages beyond available pages
		response := utils.PaginatedResponse{
			Data:       []interface{}{}, // Empty array
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}

	// Return consistent paginated response format
	response := utils.PaginatedResponse{
		Data:       brands,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
