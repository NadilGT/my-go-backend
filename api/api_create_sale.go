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

func CreateSaleApi(c *fiber.Ctx) error {
	inputObj := dto.Sale{}

	if err := c.BodyParser(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Validate the input
	validate := validator.New()
	if validationErr := validate.Struct(inputObj); validationErr != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, validationErr.Error())
	}

	ctx := context.Background()

	// Generate Sale ID
	id, err := dao.GenerateId(ctx, "Sales", "SALE")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	inputObj.SaleId = id

	// Set timestamps
	now := time.Now().UTC()
	inputObj.CreatedAt = now
	inputObj.UpdatedAt = now

	// Set sale date if not provided
	if inputObj.SaleDate.IsZero() {
		inputObj.SaleDate = now
	}

	// Calculate totals for each item and overall totals
	var subTotal float64
	for i := range inputObj.Items {
		inputObj.Items[i].TotalPrice = float64(inputObj.Items[i].Quantity) * inputObj.Items[i].UnitPrice
		subTotal += inputObj.Items[i].TotalPrice
	}
	inputObj.SubTotal = subTotal

	// Calculate grand total (SubTotal - TotalDiscount + TaxAmount)
	inputObj.GrandTotal = inputObj.SubTotal - inputObj.TotalDiscount + inputObj.TaxAmount

	// Calculate change amount
	inputObj.ChangeAmount = inputObj.PaidAmount - inputObj.GrandTotal

	// Validate payment
	if inputObj.PaidAmount < inputObj.GrandTotal {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Paid amount cannot be less than grand total")
	}

	// Set default status if not provided
	if inputObj.Status == "" {
		inputObj.Status = "completed"
	}

	// Set default deleted status
	inputObj.Deleted = false

	// Save to database
	if err := dao.DB_CreateSale(ctx, &inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Update product stock quantities
	if err := dao.DB_UpdateProductStockAfterSale(ctx, inputObj.Items); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Sale created but failed to update stock: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sale created successfully",
		"saleId":  inputObj.SaleId,
		"data":    inputObj,
	})
}
