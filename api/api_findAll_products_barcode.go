package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetProductsByBarcodeApi(c *fiber.Ctx) error {
	barcode := c.Query("barcode", "")
	if barcode == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "Barcode is required", nil)
	}

	cursor := c.Query("cursor", "")
	perPageStr := c.Query("per_page", "15")

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		perPage = 15
	}

	switch perPage {
	case 15, 25, 50:
	default:
		perPage = 15
	}

	products, nextCursor, hasMore, err := dao.DB_FindProductsByBarcodeCursorPaginated(barcode, perPage, cursor)
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
