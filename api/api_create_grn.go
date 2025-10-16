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

func CreateGRN(c *fiber.Ctx) error {
	inputObj := dto.GRN{}

	if err := c.BodyParser(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	ctx := context.Background()
	id, err := dao.GenerateId(ctx, "GRNs", "GRN")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	inputObj.GRNId = id
	now := time.Now().UTC()
	inputObj.CreatedAt = now
	inputObj.UpdatedAt = now

	// Calculate total amount for each item and overall total
	var totalAmount float64
	for i := range inputObj.Items {
		inputObj.Items[i].TotalCost = float64(inputObj.Items[i].ReceivedQty) * inputObj.Items[i].UnitCost
		totalAmount += inputObj.Items[i].TotalCost
	}
	inputObj.TotalAmount = totalAmount

	// Set default status if not provided
	if inputObj.Status == "" {
		inputObj.Status = "pending"
	}

	// Set default deleted status
	inputObj.Deleted = false

	if err := functions.UniqueCheck(inputObj, "GRNs", []string{"GRNId"}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if validationErr := validate.Struct(inputObj); validationErr != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, validationErr.Error())
	}

	err = dao.DB_CreateGRN(&inputObj)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c)
}
