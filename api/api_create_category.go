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

	// Parse request body
	if err := c.BodyParser(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Generate category ID
	ctx := context.Background()
	id, err := dao.GenerateId(ctx, "Categories", "CAT")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Assign generated values
	inputObj.CategoryId = id
	now := time.Now().UTC()
	inputObj.CreatedAt = now
	inputObj.UpdatedAt = now

	// Save to DB
	if err := dao.DB_CreateCategory(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Return success with created object
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    inputObj,
	})
}
