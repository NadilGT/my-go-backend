package api

import (
	"context"
	"employee-crud/dao"
	"employee-crud/dto"
	"employee-crud/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateCategoryApi(c *fiber.Ctx) error {
	inputObj := dto.Category{}

	categoryName := c.Query("categoryName")
	if categoryName == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "categoryName is required")
	}

	ctx := context.Background()
	id, err := dao.GenerateId(ctx, "Categories", "CAT")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	inputObj.CategoryId = id
	inputObj.Name = categoryName
	now := time.Now().UTC()
	inputObj.CreatedAt = now
	inputObj.UpdatedAt = now

	if err := dao.DB_CreateCategory(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    inputObj,
	})
}
