package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FindAllCategoriesApi(c *fiber.Ctx) error {
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
	categories, total, err := dao.DB_FindAllCategoriesPaginated(page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get product counts for all categories in this page
	categoryIds := make([]string, len(categories))
	for i, cat := range categories {
		categoryIds[i] = cat.CategoryId
	}

	if len(categoryIds) > 0 {
		productCounts, err := dao.DB_GetProductCountsForCategories(categoryIds)
		if err != nil {
			// Log error but continue without product counts
			// You can add proper logging here if needed
			productCounts = make(map[string]int64)
		}

		// Add product counts to categories
		for i := range categories {
			categories[i].ProductCount = productCounts[categories[i].CategoryId]
		}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	// Return consistent paginated response format
	response := utils.PaginatedResponse{
		Data:       categories,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
