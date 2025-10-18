package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetSaleBillApi(c *fiber.Ctx) error {
	saleId := c.Query("saleId")
	if saleId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "saleId parameter is required",
		})
	}

	sale, err := dao.DB_FindSaleById(saleId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Sale not found",
		})
	}

	billData := generateSaleBill(sale)
	return c.JSON(billData)
}

func generateSaleBill(sale *dto.Sale) fiber.Map {
	// Format dates for display
	saleDateFormatted := sale.SaleDate.Format("2006-01-02 15:04:05")
	createdAtFormatted := sale.CreatedAt.Format("2006-01-02 15:04:05")

	return fiber.Map{
		"billMetadata": fiber.Map{
			"generatedAt":  time.Now().Format("2006-01-02 15:04:05"),
			"billTitle":    "Sales Invoice",
			"billSubtitle": "Receipt for Purchase",
		},
		"saleHeader": fiber.Map{
			"saleNumber": sale.SaleNumber,
			"saleId":     sale.SaleId,

			// Customer information - only included if provided
			"customerInfo": generateCustomerInfo(sale),

			"dates": fiber.Map{
				"saleDate":           sale.SaleDate,
				"saleDateFormatted":  saleDateFormatted,
				"createdAt":          sale.CreatedAt,
				"createdAtFormatted": createdAtFormatted,
			},
			"soldBy": sale.SoldBy,
			"notes":  sale.Notes,
			"status": fiber.Map{
				"value":       sale.Status,
				"displayName": getStatusDisplayText(sale.Status),
			},
		},
		"itemsDetail": enrichItemsForBill(sale.Items),
		"summary": fiber.Map{
			"totals": fiber.Map{
				"subTotal":               sale.SubTotal,
				"totalDiscount":          sale.TotalDiscount,
				"taxAmount":              sale.TaxAmount,
				"grandTotal":             sale.GrandTotal,
				"formattedSubTotal":      formatCurrencyForSale(sale.SubTotal),
				"formattedTotalDiscount": formatCurrencyForSale(sale.TotalDiscount),
				"formattedTaxAmount":     formatCurrencyForSale(sale.TaxAmount),
				"formattedGrandTotal":    formatCurrencyForSale(sale.GrandTotal),
			},
			"payment": fiber.Map{
				"paymentMethod":         sale.PaymentMethod,
				"paidAmount":            sale.PaidAmount,
				"changeAmount":          sale.ChangeAmount,
				"formattedPaidAmount":   formatCurrencyForSale(sale.PaidAmount),
				"formattedChangeAmount": formatCurrencyForSale(sale.ChangeAmount),
			},
			"quantities": fiber.Map{
				"totalItems":    len(sale.Items),
				"totalQuantity": calculateTotalQuantity(sale.Items),
			},
		},
	}
}

func generateCustomerInfo(sale *dto.Sale) fiber.Map {
	// Only include customer info if at least name is provided
	if sale.CustomerName != "" {
		return fiber.Map{
			"hasCustomer": true,
			"name":        sale.CustomerName,
			"phone":       sale.CustomerPhone, // Can be empty
		}
	}

	return fiber.Map{
		"hasCustomer": false,
	}
}

func enrichItemsForBill(items []dto.SaleItem) []fiber.Map {
	enrichedItems := make([]fiber.Map, len(items))

	for i, item := range items {
		enrichedItems[i] = fiber.Map{
			"productId":   item.ProductId,
			"productName": item.ProductName,
			"quantity":    item.Quantity,
			"unitPrice":   item.UnitPrice,
			"totalPrice":  item.TotalPrice,
			"discount":    item.Discount,
			"formatted": fiber.Map{
				"unitPrice":  formatCurrencyForSale(item.UnitPrice),
				"totalPrice": formatCurrencyForSale(item.TotalPrice),
				"discount":   formatCurrencyForSale(item.Discount),
			},
		}
	}

	return enrichedItems
}

func getStatusDisplayText(status string) string {
	statusMap := map[string]string{
		"completed":      "Completed",
		"refunded":       "Refunded",
		"partial_refund": "Partially Refunded",
	}

	if displayName, exists := statusMap[status]; exists {
		return displayName
	}
	return status
}

func calculateTotalQuantity(items []dto.SaleItem) int {
	total := 0
	for _, item := range items {
		total += item.Quantity
	}
	return total
}

func formatCurrencyForSale(amount float64) string {
	return "Rs. " + strconv.FormatFloat(amount, 'f', 2, 64)
}
