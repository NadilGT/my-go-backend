package api

import (
	"employee-crud/dto"

	"github.com/gofiber/fiber/v2"
)

func CalculateOrderSummaryApi(c *fiber.Ctx) error {
	var req dto.CalculateOrderSummaryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate that items array is not empty
	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one item is required",
		})
	}

	// Calculate subtotal
	var subtotal float64 = 0
	for _, item := range req.Items {
		subtotal += item.TotalPrice
	}

	// Calculate tax
	var tax float64 = 0
	if req.TaxPercentage > 0 {
		tax = subtotal * (req.TaxPercentage / 100)
	}

	// Calculate discount
	var discount float64 = 0
	if req.DiscountType == "percentage" {
		discount = subtotal * (req.Discount / 100)
	} else {
		discount = req.Discount
	}

	// Calculate total
	total := subtotal + tax - discount

	response := dto.OrderSummaryResponse{
		Subtotal: subtotal,
		Tax:      tax,
		Discount: discount,
		Total:    total,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
