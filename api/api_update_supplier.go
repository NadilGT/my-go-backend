package api

import (
	"context"
	"employee-crud/dao"
	"employee-crud/dto"
	"employee-crud/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UpdateSupplierApi(c *fiber.Ctx) error {
	inputObj := dto.Supplier{}

	if err := c.BodyParser(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if inputObj.SupplierId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SupplierId is required")
	}

	validate := validator.New()
	if validationErr := validate.StructPartial(inputObj,
		"Name", "Contact", "Email", "Address"); validationErr != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, validationErr.Error())
	}

	inputObj.UpdatedAt = time.Now().UTC()

	if err := dao.DB_UpdateSupplier(context.Background(), &inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c)
}
