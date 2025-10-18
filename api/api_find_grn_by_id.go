package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FindGRNByIdApi(c *fiber.Ctx) error {
	grnId := c.Query("grnId")
	if grnId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "grnId parameter is required",
		})
	}

	grn, err := dao.DB_FindGRNById(grnId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "GRN not found",
		})
	}

	// Format the response with report-ready structure
	response := fiber.Map{
		"grn": grn,
		"reportData": fiber.Map{
			"grnDetails": fiber.Map{
				"grnNumber":     grn.GRNNumber,
				"grnId":         grn.GRNId,
				"supplierName":  grn.SupplierName,
				"supplierId":    grn.SupplierId,
				"receivedDate":  grn.ReceivedDate,
				"invoiceNumber": grn.InvoiceNumber,
				"invoiceDate":   grn.InvoiceDate,
				"status":        grn.Status,
				"receivedBy":    grn.ReceivedBy,
				"notes":         grn.Notes,
				"totalAmount":   grn.TotalAmount,
			},
			"items": grn.Items,
			"summary": fiber.Map{
				"totalItems":       len(grn.Items),
				"totalExpectedQty": calculateTotalExpectedQty(grn.Items),
				"totalReceivedQty": calculateTotalReceivedQty(grn.Items),
				"totalAmount":      grn.TotalAmount,
				"formattedAmount":  formatCurrency(grn.TotalAmount),
			},
		},
	}

	return c.JSON(response)
}

func calculateTotalExpectedQty(items []dto.GRNItem) int {
	total := 0
	for _, item := range items {
		total += item.ExpectedQty
	}
	return total
}

func calculateTotalReceivedQty(items []dto.GRNItem) int {
	total := 0
	for _, item := range items {
		total += item.ReceivedQty
	}
	return total
}

func formatCurrency(amount float64) string {
	return "Rs. " + strconv.FormatFloat(amount, 'f', 2, 64)
}
