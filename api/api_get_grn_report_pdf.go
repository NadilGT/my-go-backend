package api

import (
	"bytes"
	"employee-crud/dao"
	"employee-crud/dto"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
)

func GetGRNReportPDFApi(c *fiber.Ctx) error {
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

	// Generate PDF
	pdfBytes, err := generateGRNReportPDF(grn)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate PDF: " + err.Error(),
		})
	}

	// Set appropriate headers for PDF download
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=GRN-Report-%s-%s.pdf",
		grn.GRNNumber, time.Now().Format("2006-01-02")))
	c.Set("Content-Length", strconv.Itoa(len(pdfBytes)))

	return c.Send(pdfBytes)
}

func generateGRNReportPDF(grn *dto.GRN) ([]byte, error) {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set margins
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)

	// Title and header
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 15, "Goods Receipt Note Report", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 8, "Generated on: "+time.Now().Format("2006-01-02 15:04:05"), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// GRN Header Information
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 10, "GRN Details", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Create two-column layout for GRN details
	pdf.SetFont("Arial", "", 10)

	// Left column
	leftColX := 15.0
	rightColX := 105.0
	currentY := pdf.GetY()

	// Left column data
	pdf.SetXY(leftColX, currentY)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "GRN Number:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, grn.GRNNumber)
	pdf.Ln(8)

	pdf.SetX(leftColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "GRN ID:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, grn.GRNId)
	pdf.Ln(8)

	pdf.SetX(leftColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Status:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, getStatusDisplayName(grn.Status))
	pdf.Ln(8)

	pdf.SetX(leftColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Received By:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, grn.ReceivedBy)
	pdf.Ln(8)

	// Right column data
	pdf.SetXY(rightColX, currentY)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Supplier:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, grn.SupplierName)
	pdf.Ln(8)

	pdf.SetX(rightColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Supplier ID:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, grn.SupplierId)
	pdf.Ln(8)

	pdf.SetX(rightColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Invoice Number:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, grn.InvoiceNumber)
	pdf.Ln(8)

	pdf.SetX(rightColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Received Date:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, grn.ReceivedDate.Format("2006-01-02 15:04:05"))
	pdf.Ln(15)

	// Summary section
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "Summary", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Calculate summary data
	totalExpected := 0
	totalReceived := 0
	itemsWithDiscrepancy := 0

	for _, item := range grn.Items {
		totalExpected += item.ExpectedQty
		totalReceived += item.ReceivedQty
		if item.ExpectedQty != item.ReceivedQty {
			itemsWithDiscrepancy++
		}
	}

	completionRate := float64(0)
	if totalExpected > 0 {
		completionRate = float64(totalReceived) / float64(totalExpected) * 100
		if completionRate > 100 {
			completionRate = 100
		}
	}

	pdf.SetFont("Arial", "", 10)
	currentY = pdf.GetY()

	// Summary left column
	pdf.SetXY(leftColX, currentY)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Total Items:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, strconv.Itoa(len(grn.Items)))
	pdf.Ln(8)

	pdf.SetX(leftColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Expected Qty:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, strconv.Itoa(totalExpected))
	pdf.Ln(8)

	pdf.SetX(leftColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Received Qty:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, strconv.Itoa(totalReceived))
	pdf.Ln(8)

	// Summary right column
	pdf.SetXY(rightColX, currentY)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Total Amount:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "Rs. "+strconv.FormatFloat(grn.TotalAmount, 'f', 2, 64))
	pdf.Ln(8)

	pdf.SetX(rightColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Discrepancies:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, strconv.Itoa(itemsWithDiscrepancy))
	pdf.Ln(8)

	pdf.SetX(rightColX)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Completion Rate:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, strconv.FormatFloat(completionRate, 'f', 1, 64)+"%")
	pdf.Ln(15)

	// Items table
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "Items Detail", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table header
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(240, 240, 240)

	colWidths := []float64{50, 25, 25, 25, 30, 35}
	headers := []string{"Product", "Expected", "Received", "Status", "Unit Cost", "Total Cost"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(8)

	// Table rows
	pdf.SetFont("Arial", "", 7)
	pdf.SetFillColor(255, 255, 255)

	for _, item := range grn.Items {
		// Check if we need a new page
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}

		discrepancy := item.ExpectedQty - item.ReceivedQty
		status := "Exact"
		if discrepancy > 0 {
			status = "Short"
		} else if discrepancy < 0 {
			status = "Excess"
		}

		// Product name (with wrapping if needed)
		productName := item.ProductName
		if len(productName) > 30 {
			productName = productName[:27] + "..."
		}

		rowData := []string{
			productName,
			strconv.Itoa(item.ExpectedQty),
			strconv.Itoa(item.ReceivedQty),
			status,
			"Rs. " + strconv.FormatFloat(item.UnitCost, 'f', 2, 64),
			"Rs. " + strconv.FormatFloat(item.TotalCost, 'f', 2, 64),
		}

		for i, data := range rowData {
			align := "L"
			if i > 0 { // Numbers and amounts should be right-aligned
				align = "R"
			}
			if i == 3 { // Status should be centered
				align = "C"
			}
			pdf.CellFormat(colWidths[i], 8, data, "1", 0, align, false, 0, "")
		}
		pdf.Ln(8)
	}

	// Notes section if exists
	if grn.Notes != "" {
		pdf.Ln(10)
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 8, "Notes:", "", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 6, grn.Notes, "", "L", false)
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 6, "This is a system generated document.", "", 1, "C", false, 0, "")

	// Get PDF bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getStatusDisplayName(status string) string {
	switch status {
	case "pending":
		return "Pending"
	case "completed":
		return "Completed"
	case "partial_received":
		return "Partially Received"
	default:
		return "Unknown"
	}
}
