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

func GetDailySalesSummaryPDFApi(c *fiber.Ctx) error {
	// Get date parameter from query (format: YYYY-MM-DD)
	// If not provided, use today's date
	dateStr := c.Query("date")
	var targetDate time.Time
	var err error

	if dateStr == "" {
		// Use today's date in Sri Lanka timezone (UTC+5:30)
		sriLankaLoc := time.FixedZone("Asia/Colombo", 5*3600+30*60)
		targetDate = time.Now().In(sriLankaLoc)
	} else {
		// Parse the provided date
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Use YYYY-MM-DD (e.g., 2025-10-19)",
			})
		}
	}

	// Get sales summary for the date
	summary, err := dao.GetDailySalesSummary(targetDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve sales summary: " + err.Error(),
		})
	}

	// Generate PDF
	pdfBytes, err := generateDailySalesSummaryPDF(summary)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate PDF: " + err.Error(),
		})
	}

	// Set appropriate headers for PDF download
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=Daily-Sales-Report-%s.pdf",
		targetDate.Format("2006-01-02")))
	c.Set("Content-Length", strconv.Itoa(len(pdfBytes)))

	return c.Send(pdfBytes)
}

func generateDailySalesSummaryPDF(summary *dto.DailySalesSummary) ([]byte, error) {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set margins
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)

	// Title and header
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 15, "Daily Sales Summary Report", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(60, 60, 60)
	pdf.CellFormat(0, 8, "Report Date: "+summary.ReportDate.Format("Monday, January 2, 2006"), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 6, "Generated on: "+time.Now().Format("2006-01-02 15:04:05"), "", 1, "C", false, 0, "")
	pdf.Ln(12)

	// Sales Overview Section
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 10, "Sales Overview", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Create bordered box for overview
	currentY := pdf.GetY()
	pdf.SetFillColor(245, 248, 250)
	pdf.Rect(15, currentY, 180, 52, "FD")

	pdf.SetXY(15, currentY+5)

	// Layout the overview in a grid (2x3)
	leftColX := 20.0
	rightColX := 105.0
	rowHeight := 15.0

	// Row 1
	currentY = pdf.GetY()

	// Total Sales
	pdf.SetXY(leftColX, currentY)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(40, 6, "Total Sales:")
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 102, 204)
	pdf.Cell(0, 6, strconv.Itoa(summary.TotalSales))

	// Total Revenue
	pdf.SetXY(rightColX, currentY)
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(40, 6, "Total Revenue:")
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 153, 51)
	pdf.Cell(0, 6, "Rs. "+strconv.FormatFloat(summary.TotalRevenue, 'f', 2, 64))

	// Row 2
	currentY += rowHeight
	pdf.SetTextColor(0, 0, 0)

	// Cash Sales
	pdf.SetXY(leftColX, currentY)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 6, "Cash Sales:")
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 6, strconv.Itoa(summary.CashSales)+" (Rs. "+strconv.FormatFloat(summary.CashRevenue, 'f', 2, 64)+")")

	// Card Sales
	pdf.SetXY(rightColX, currentY)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 6, "Card Sales:")
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 6, strconv.Itoa(summary.CardSales)+" (Rs. "+strconv.FormatFloat(summary.CardRevenue, 'f', 2, 64)+")")

	// Row 3
	currentY += rowHeight

	// Total Discount
	pdf.SetXY(leftColX, currentY)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 6, "Total Discount:")
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(204, 0, 0)
	pdf.Cell(0, 6, "Rs. "+strconv.FormatFloat(summary.TotalDiscount, 'f', 2, 64))

	// Total Tax
	pdf.SetXY(rightColX, currentY)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(40, 6, "Total Tax:")
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 6, "Rs. "+strconv.FormatFloat(summary.TotalTax, 'f', 2, 64))

	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(20)

	// Top Selling Items Section
	if len(summary.TopSellingItems) > 0 {
		pdf.SetFont("Arial", "B", 16)
		pdf.CellFormat(0, 10, "Top Selling Items", "", 1, "L", false, 0, "")
		pdf.Ln(2)

		// Table header
		pdf.SetFont("Arial", "B", 9)
		pdf.SetFillColor(52, 73, 94)
		pdf.SetTextColor(255, 255, 255)

		colWidths := []float64{20, 70, 30, 30, 30}
		headers := []string{"Rank", "Product Name", "Qty Sold", "Unit Price", "Total Amount"}

		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(8)

		// Table rows
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(0, 0, 0)

		for idx, item := range summary.TopSellingItems {
			// Check if we need a new page
			if pdf.GetY() > 250 {
				pdf.AddPage()
			}

			// Alternate row colors
			if idx%2 == 0 {
				pdf.SetFillColor(245, 245, 245)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}

			// Product name (with wrapping if needed)
			productName := item.ProductName
			if len(productName) > 40 {
				productName = productName[:37] + "..."
			}

			rank := strconv.Itoa(idx + 1)
			rowData := []string{
				rank,
				productName,
				strconv.Itoa(item.Quantity),
				"Rs. " + strconv.FormatFloat(item.UnitPrice, 'f', 2, 64),
				"Rs. " + strconv.FormatFloat(item.TotalAmount, 'f', 2, 64),
			}

			for i, data := range rowData {
				align := "L"
				if i == 0 { // Rank centered
					align = "C"
				}
				if i > 1 { // Numbers should be right-aligned
					align = "R"
				}
				pdf.CellFormat(colWidths[i], 8, data, "1", 0, align, true, 0, "")
			}
			pdf.Ln(8)
		}

		pdf.Ln(5)
	}

	// All Products Sold Section
	if len(summary.ProductsSold) > 0 {
		pdf.SetFont("Arial", "B", 16)
		pdf.CellFormat(0, 10, "All Products Sold", "", 1, "L", false, 0, "")
		pdf.Ln(2)

		pdf.SetFont("Arial", "", 9)
		pdf.SetTextColor(100, 100, 100)
		pdf.CellFormat(0, 5, fmt.Sprintf("Total Products: %d", len(summary.ProductsSold)), "", 1, "L", false, 0, "")
		pdf.Ln(3)

		// Table header
		pdf.SetFont("Arial", "B", 9)
		pdf.SetFillColor(52, 73, 94)
		pdf.SetTextColor(255, 255, 255)

		colWidths := []float64{25, 70, 30, 30, 25}
		headers := []string{"Product ID", "Product Name", "Qty Sold", "Unit Price", "Total"}

		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(8)

		// Table rows
		pdf.SetFont("Arial", "", 7)
		pdf.SetTextColor(0, 0, 0)

		for idx, item := range summary.ProductsSold {
			// Check if we need a new page
			if pdf.GetY() > 260 {
				pdf.AddPage()

				// Repeat header on new page
				pdf.SetFont("Arial", "B", 9)
				pdf.SetFillColor(52, 73, 94)
				pdf.SetTextColor(255, 255, 255)
				for i, header := range headers {
					pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
				}
				pdf.Ln(8)
				pdf.SetFont("Arial", "", 7)
				pdf.SetTextColor(0, 0, 0)
			}

			// Alternate row colors
			if idx%2 == 0 {
				pdf.SetFillColor(250, 250, 250)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}

			// Product name (with wrapping if needed)
			productName := item.ProductName
			if len(productName) > 35 {
				productName = productName[:32] + "..."
			}

			rowData := []string{
				item.ProductID,
				productName,
				strconv.Itoa(item.Quantity),
				"Rs. " + strconv.FormatFloat(item.UnitPrice, 'f', 2, 64),
				"Rs. " + strconv.FormatFloat(item.TotalAmount, 'f', 2, 64),
			}

			for i, data := range rowData {
				align := "L"
				if i > 1 { // Numbers should be right-aligned
					align = "R"
				}
				pdf.CellFormat(colWidths[i], 7, data, "1", 0, align, true, 0, "")
			}
			pdf.Ln(7)
		}
	}

	// Summary footer box
	pdf.Ln(10)
	currentY = pdf.GetY()

	// Check if we need a new page for footer
	if currentY > 250 {
		pdf.AddPage()
		currentY = pdf.GetY()
	}

	pdf.SetFillColor(52, 73, 94)
	pdf.SetTextColor(255, 255, 255)
	pdf.Rect(15, currentY, 180, 20, "FD")

	pdf.SetXY(15, currentY+5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(90, 10, "Grand Total Revenue:", "", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(90, 10, "Rs. "+strconv.FormatFloat(summary.TotalRevenue, 'f', 2, 64), "", 0, "L", false, 0, "")

	// Footer
	pdf.SetY(-25)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 5, "This is a system generated document.", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Report covers sales from %s 00:00:00 to %s 23:59:59 (UTC)",
		summary.ReportDate.Format("2006-01-02"),
		summary.ReportDate.Format("2006-01-02")), "", 1, "C", false, 0, "")

	// Get PDF bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
