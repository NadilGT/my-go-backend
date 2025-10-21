package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateSaleApi(c *fiber.Ctx) error {
	var req dto.CreateSaleRequest

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

	// Validate payment method
	if req.PaymentMethod != "cash" && req.PaymentMethod != "card" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Payment method must be either 'cash' or 'card'",
		})
	}

	// Calculate order summary
	var subtotal float64 = 0
	for i := range req.Items {
		// Verify product exists and has sufficient stock
		product, err := dao.GetProductByProductId(req.Items[i].ProductID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found: " + req.Items[i].ProductID,
			})
		}

		if product.StockQty < req.Items[i].Quantity {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Insufficient stock for product: " + product.Name,
			})
		}

		// Calculate total price for this item
		req.Items[i].ProductName = product.Name
		req.Items[i].UnitPrice = product.SellingPrice
		req.Items[i].TotalPrice = product.SellingPrice * float64(req.Items[i].Quantity)
		subtotal += req.Items[i].TotalPrice
	}

	// Calculate tax
	var tax float64 = 0
	if req.TaxPercentage > 0 {
		tax = subtotal * (req.TaxPercentage / 100)
	} else {
		tax = req.Tax
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

	// Calculate change for cash payments
	var change float64 = 0
	if req.PaymentMethod == "cash" {
		if req.AmountReceived < total {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Amount received is less than total",
			})
		}
		change = req.AmountReceived - total
	}

	// Create sale object
	sale := &dto.Sale{
		SaleID:         uuid.New().String(),
		CustomerName:   req.CustomerName,
		MobileNumber:   req.MobileNumber,
		Items:          req.Items,
		Subtotal:       subtotal,
		Tax:            tax,
		TaxPercentage:  req.TaxPercentage,
		Discount:       discount,
		DiscountType:   req.DiscountType,
		Total:          total,
		PaymentMethod:  req.PaymentMethod,
		AmountReceived: req.AmountReceived,
		Change:         change,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save sale to database
	if err := dao.CreateSale(sale); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create sale",
		})
	}

	// Update product stock for each item
	// This automatically syncs each product to the Stocks collection
	for _, item := range req.Items {
		if err := dao.UpdateProductStock(item.ProductID, item.Quantity); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to update stock for product: " + item.ProductID,
				"details": err.Error(),
			})
		}
	}

	// Return success
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sale created successfully and stocks updated",
		"sale":    sale,
	})
}
