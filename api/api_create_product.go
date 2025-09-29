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

func CreateProduct(c *fiber.Ctx) error {
	inputObj := dto.Product{}

	if err := c.BodyParser(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	ctx := context.Background()
	id, err := dao.GenerateId(ctx, "Products", "PRD")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	inputObj.ProductID = id
	now := time.Now().UTC()
	inputObj.CreatedAt = now
	inputObj.UpdatedAt = now

	if err := functions.UniqueCheck(inputObj, "Products", []string{"ProductId"}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if validationErr := validate.Struct(inputObj); validationErr != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, validationErr.Error())
	}
	err = dao.DB_CreateProduct(&inputObj)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c)
}
