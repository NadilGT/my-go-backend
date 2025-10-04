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

func UpdateProductApi(c *fiber.Ctx) error {
	inputObj := dto.Product{}

	if err := c.BodyParser(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if inputObj.ProductId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "ProductId is required")
	}

	validate := validator.New()
	if validationErr := validate.StructPartial(inputObj,
		"Name", "Barcode", "CategoryID", "BrandID", "SubCategoryID",
		"CostPrice", "SellingPrice", "StockQty", "ExpiryDate", "Deleted"); validationErr != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, validationErr.Error())
	}

	inputObj.UpdatedAt = time.Now().UTC()

	if err := dao.DB_UpdateProduct(context.Background(), &inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c)
}
