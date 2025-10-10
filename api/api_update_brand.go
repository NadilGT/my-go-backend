package api

import (
	"context"
	"employee-crud/dao"
	"employee-crud/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UpdateBrandApi(c *fiber.Ctx) error {
	brandId := c.Query("brandId")
	name := c.Query("name")
	categoryId := c.Query("categoryId")

	if brandId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BrandId is required")
	}
	if name == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Name is required")
	}
	if categoryId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "CategoryId is required")
	}

	updatedAt := time.Now().UTC()

	if err := dao.DB_UpdateBrand(context.Background(), brandId, name, categoryId, updatedAt); err != nil {
		if err.Error() == "not_found" {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, "Brand not found")
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c)
}
