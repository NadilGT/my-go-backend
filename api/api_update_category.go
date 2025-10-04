package api

import (
	"context"
	"employee-crud/dao"
	"employee-crud/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UpdateCategoryApi(c *fiber.Ctx) error {
	categoryId := c.Query("categoryId")
	name := c.Query("name")

	if categoryId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "CategoryId is required")
	}
	if name == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Name is required")
	}

	updatedAt := time.Now().UTC()

	if err := dao.DB_UpdateCategory(context.Background(), categoryId, name, updatedAt); err != nil {
		if err.Error() == "not_found" {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, "Category not found")
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c)
}
