package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetGRNReportApi(c *fiber.Ctx) error {
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

	// Enhanced report structure with all necessary data for frontend
	reportData := generateGRNReport(grn)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    reportData,
	})
}

func generateGRNReport(grn *dto.GRN) fiber.Map {
	// Calculate summary statistics
	totalExpected := 0
	totalReceived := 0
	itemsWithDiscrepancy := 0
	discrepancyPercentage := 0.0

	for _, item := range grn.Items {
		totalExpected += item.ExpectedQty
		totalReceived += item.ReceivedQty
		if item.ExpectedQty != item.ReceivedQty {
			itemsWithDiscrepancy++
		}
	}

	if totalExpected > 0 {
		discrepancyPercentage = float64(totalExpected-totalReceived) / float64(totalExpected) * 100
	}

	// Format dates for display
	receivedDateFormatted := grn.ReceivedDate.Format("2006-01-02 15:04:05")
	invoiceDateFormatted := ""
	if grn.InvoiceDate != nil {
		invoiceDateFormatted = grn.InvoiceDate.Format("2006-01-02")
	}

	// Status styling information
	statusInfo := getStatusInfo(grn.Status)

	return fiber.Map{
		"reportMetadata": fiber.Map{
			"generatedAt":    time.Now().Format("2006-01-02 15:04:05"),
			"reportTitle":    "Goods Receipt Note Report",
			"reportSubtitle": "GRN Details and Item Summary",
		},
		"grnHeader": fiber.Map{
			"grnNumber": grn.GRNNumber,
			"grnId":     grn.GRNId,
			"supplierInfo": fiber.Map{
				"id":   grn.SupplierId,
				"name": grn.SupplierName,
			},
			"dates": fiber.Map{
				"received":          grn.ReceivedDate,
				"receivedFormatted": receivedDateFormatted,
				"invoice":           grn.InvoiceDate,
				"invoiceFormatted":  invoiceDateFormatted,
			},
			"invoiceNumber": grn.InvoiceNumber,
			"receivedBy":    grn.ReceivedBy,
			"notes":         grn.Notes,
			"status": fiber.Map{
				"value":       grn.Status,
				"displayName": statusInfo["displayName"],
				"color":       statusInfo["color"],
				"badge":       statusInfo["badge"],
			},
		},
		"itemsDetail": enrichItemsForReport(grn.Items),
		"summary": fiber.Map{
			"financials": fiber.Map{
				"totalAmount":          grn.TotalAmount,
				"formattedTotalAmount": "Rs. " + strconv.FormatFloat(grn.TotalAmount, 'f', 2, 64),
				"currency":             "LKR",
			},
			"quantities": fiber.Map{
				"totalItems":            len(grn.Items),
				"totalExpectedQty":      totalExpected,
				"totalReceivedQty":      totalReceived,
				"itemsWithDiscrepancy":  itemsWithDiscrepancy,
				"discrepancyPercentage": discrepancyPercentage,
			},
			"completion": fiber.Map{
				"isComplete":     grn.Status == "completed",
				"isPartial":      grn.Status == "partial_received",
				"isPending":      grn.Status == "pending",
				"completionRate": calculateCompletionRate(grn.Items),
			},
		},
		"printOptions": fiber.Map{
			"includeHeader":    true,
			"includeFooter":    true,
			"includeSignature": true,
			"paperSize":        "A4",
			"orientation":      "portrait",
		},
	}
}

func enrichItemsForReport(items []dto.GRNItem) []fiber.Map {
	enrichedItems := make([]fiber.Map, len(items))

	for i, item := range items {
		discrepancy := item.ExpectedQty - item.ReceivedQty
		discrepancyStatus := "exact"
		if discrepancy > 0 {
			discrepancyStatus = "shortage"
		} else if discrepancy < 0 {
			discrepancyStatus = "excess"
		}

		expiryFormatted := ""
		if item.ExpiryDate != nil {
			expiryFormatted = item.ExpiryDate.Format("2006-01-02")
		}

		enrichedItems[i] = fiber.Map{
			"productId":   item.ProductId,
			"productName": item.ProductName,
			"quantities": fiber.Map{
				"expected":    item.ExpectedQty,
				"received":    item.ReceivedQty,
				"discrepancy": discrepancy,
				"status":      discrepancyStatus,
			},
			"costs": fiber.Map{
				"unitCost":           item.UnitCost,
				"totalCost":          item.TotalCost,
				"formattedUnitCost":  "Rs. " + strconv.FormatFloat(item.UnitCost, 'f', 2, 64),
				"formattedTotalCost": "Rs. " + strconv.FormatFloat(item.TotalCost, 'f', 2, 64),
			},
			"additionalInfo": fiber.Map{
				"expiryDate":      item.ExpiryDate,
				"expiryFormatted": expiryFormatted,
				"batchNumber":     item.BatchNumber,
				"remarks":         item.Remarks,
			},
		}
	}

	return enrichedItems
}

func getStatusInfo(status string) fiber.Map {
	statusMap := map[string]fiber.Map{
		"pending": {
			"displayName": "Pending",
			"color":       "#FFA500",
			"badge":       "warning",
		},
		"completed": {
			"displayName": "Completed",
			"color":       "#28A745",
			"badge":       "success",
		},
		"partial_received": {
			"displayName": "Partially Received",
			"color":       "#17A2B8",
			"badge":       "info",
		},
	}

	if info, exists := statusMap[status]; exists {
		return info
	}

	return fiber.Map{
		"displayName": "Unknown",
		"color":       "#6C757D",
		"badge":       "secondary",
	}
}

func calculateCompletionRate(items []dto.GRNItem) float64 {
	if len(items) == 0 {
		return 0
	}

	totalExpected := 0
	totalReceived := 0

	for _, item := range items {
		totalExpected += item.ExpectedQty
		totalReceived += item.ReceivedQty
	}

	if totalExpected == 0 {
		return 0
	}

	rate := float64(totalReceived) / float64(totalExpected) * 100
	if rate > 100 {
		return 100
	}
	return rate
}
