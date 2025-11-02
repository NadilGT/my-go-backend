package api

import (
	"context"
	"employee-crud/dao"
	"employee-crud/dto"
	"employee-crud/functions"
	"employee-crud/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateSupplier(c *fiber.Ctx) error {
	inputObj := dto.Supplier{}

	if err := c.BodyParser(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	ctx := context.Background()
	id, err := dao.GenerateId(ctx, "Suppliers", "SUPl")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	inputObj.SupplierId = id
	now := time.Now().UTC()
	inputObj.CreatedAt = now
	inputObj.UpdatedAt = now

	// Set status to active by default
	inputObj.Status = "active"

	if err := functions.UniqueCheck(inputObj, "Suppliers", []string{"SupplierId"}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if validationErr := validate.Struct(inputObj); validationErr != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, validationErr.Error())
	}
	err = dao.DB_CreateSupplier(&inputObj)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c)
}
