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

func CreateBrand(c *fiber.Ctx) error {
	inputObj := dto.Brand{}

	if err := c.BodyParser(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	ctx := context.Background()
	id, err := dao.GenerateId(ctx, "Brands", "BRD")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	inputObj.BrandId = id
	now := time.Now().UTC()
	inputObj.CreatedAt = now
	inputObj.UpdatedAt = now

	if err := functions.UniqueCheck(inputObj, "Brands", []string{"BrandId"}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if validationErr := validate.Struct(inputObj); validationErr != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, validationErr.Error())
	}
	err = dao.DB_CreateBrand(&inputObj)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c)
}
