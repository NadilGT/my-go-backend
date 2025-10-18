package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetSalesLast24HoursApi returns all sales from the last 24 hours
func GetSalesLast24HoursApi(c *fiber.Ctx) error {
	sales, err := dao.DB_FindSalesLast24Hours()
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Calculate summary statistics
	summary := calculateSalesSummary(sales)

	return c.JSON(fiber.Map{
		"sales":   sales,
		"summary": summary,
		"period": fiber.Map{
			"from": time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
			"to":   time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

func calculateSalesSummary(sales []interface{}) fiber.Map {
	totalSales := len(sales)
	totalRevenue := 0.0
	totalDiscount := 0.0
	totalTax := 0.0
	totalItems := 0

	cashSales := 0
	cardSales := 0
	transferSales := 0
	cashRevenue := 0.0
	cardRevenue := 0.0
	transferRevenue := 0.0

	salesWithCustomer := 0
	salesWithoutCustomer := 0

	for _, s := range sales {
		if saleMap, ok := s.(map[string]interface{}); ok {
			// Revenue calculation
			if grandTotal, ok := saleMap["grandTotal"].(float64); ok {
				totalRevenue += grandTotal
			}

			// Discount calculation
			if discount, ok := saleMap["totalDiscount"].(float64); ok {
				totalDiscount += discount
			}

			// Tax calculation
			if tax, ok := saleMap["taxAmount"].(float64); ok {
				totalTax += tax
			}

			// Items count
			if items, ok := saleMap["items"].([]interface{}); ok {
				for _, item := range items {
					if itemMap, ok := item.(map[string]interface{}); ok {
						if qty, ok := itemMap["quantity"].(float64); ok {
							totalItems += int(qty)
						} else if qty, ok := itemMap["quantity"].(int); ok {
							totalItems += qty
						}
					}
				}
			}

			// Payment method breakdown
			if paymentMethod, ok := saleMap["paymentMethod"].(string); ok {
				grandTotal := 0.0
				if gt, ok := saleMap["grandTotal"].(float64); ok {
					grandTotal = gt
				}

				switch paymentMethod {
				case "cash":
					cashSales++
					cashRevenue += grandTotal
				case "card":
					cardSales++
					cardRevenue += grandTotal
				case "transfer":
					transferSales++
					transferRevenue += grandTotal
				}
			}

			// Customer info tracking
			if customerName, ok := saleMap["customerName"].(string); ok && customerName != "" {
				salesWithCustomer++
			} else {
				salesWithoutCustomer++
			}
		}
	}

	return fiber.Map{
		"totalSales":     totalSales,
		"totalRevenue":   totalRevenue,
		"totalDiscount":  totalDiscount,
		"totalTax":       totalTax,
		"totalItemsSold": totalItems,
		"paymentMethods": fiber.Map{
			"cash": fiber.Map{
				"count":   cashSales,
				"revenue": cashRevenue,
			},
			"card": fiber.Map{
				"count":   cardSales,
				"revenue": cardRevenue,
			},
			"transfer": fiber.Map{
				"count":   transferSales,
				"revenue": transferRevenue,
			},
		},
		"customerInfo": fiber.Map{
			"withCustomer":    salesWithCustomer,
			"withoutCustomer": salesWithoutCustomer,
		},
	}
}
