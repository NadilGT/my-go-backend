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

// GetDateRangeReportsPDFApi downloads all saved reports from selected date to end of month
func GetDateRangeReportsPDFApi(c *fiber.Ctx) error {
	// Get start date parameter from query (format: YYYY-MM-DD)
	dateStr := c.Query("startDate")
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "startDate parameter is required. Use format YYYY-MM-DD (e.g., 2025-10-15)",
		})
	}

	// Parse the provided date
	startDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Use YYYY-MM-DD (e.g., 2025-10-15)",
		})
	}

	// Convert to Sri Lanka timezone
	sriLankaLoc := time.FixedZone("Asia/Colombo", 5*3600+30*60)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, sriLankaLoc)

	// Calculate end of month
	year := startDate.Year()
	month := startDate.Month()
	endOfMonth := time.Date(year, month+1, 0, 23, 59, 59, 0, sriLankaLoc) // Last day of the month

	// Check if the requested date is in the future
	now := time.Now().In(sriLankaLoc)
	if startDate.After(now) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot generate reports for future dates",
		})
	}

	// If end of month is in the future, use yesterday as the end date
	if endOfMonth.After(now) {
		endOfMonth = now.AddDate(0, 0, -1) // Yesterday
		endOfMonth = time.Date(endOfMonth.Year(), endOfMonth.Month(), endOfMonth.Day(), 23, 59, 59, 0, sriLankaLoc)
	}

	// Fetch all reports for the month
	reports, err := dao.GetDailyReportsByMonth(year, int(month))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve reports: " + err.Error(),
		})
	}

	// Filter reports from startDate to endOfMonth
	var filteredReports []dto.DailyReportDocument
	for _, report := range reports {
		reportDate := report.ReportDate.In(sriLankaLoc)
		if (reportDate.Equal(startDate) || reportDate.After(startDate)) &&
			(reportDate.Equal(endOfMonth) || reportDate.Before(endOfMonth)) {
			filteredReports = append(filteredReports, report)
		}
	}

	if len(filteredReports) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("No saved reports found from %s to %s. Reports may not be generated yet or have been deleted.",
				startDate.Format("2006-01-02"), endOfMonth.Format("2006-01-02")),
		})
	}

	// Generate combined PDF
	pdfBytes, err := generateDateRangeReportsPDF(filteredReports, startDate, endOfMonth)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate PDF: " + err.Error(),
		})
	}

	// Set appropriate headers for PDF download
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=Sales-Reports-%s-to-%s.pdf",
		startDate.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02")))
	c.Set("Content-Length", strconv.Itoa(len(pdfBytes)))

	return c.Send(pdfBytes)
}

func generateDateRangeReportsPDF(reports []dto.DailyReportDocument, startDate, endDate time.Time) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)

	// Add cover page
	pdf.AddPage()
	addCoverPage(pdf, reports, startDate, endDate)

	// Add each daily report
	for idx, report := range reports {
		pdf.AddPage()
		addDailyReportPage(pdf, &report, idx+1, len(reports))
	}

	// Get PDF bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func addCoverPage(pdf *gofpdf.Fpdf, reports []dto.DailyReportDocument, startDate, endDate time.Time) {
	// Title
	pdf.SetFont("Arial", "B", 28)
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(30)
	pdf.CellFormat(0, 15, "Sales Reports Summary", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(60, 60, 60)
	pdf.CellFormat(0, 10, fmt.Sprintf("From %s to %s",
		startDate.Format("January 2, 2006"),
		endDate.Format("January 2, 2006")), "", 1, "C", false, 0, "")

	pdf.Ln(20)

	// Calculate totals
	var totalRevenue, totalDiscount, totalTax float64
	var totalSalesCount int
	for _, report := range reports {
		totalRevenue += report.TotalRevenue
		totalDiscount += report.TotalDiscount
		totalTax += report.TotalTax
		totalSalesCount += report.TotalSales
	}

	// Summary box
	currentY := pdf.GetY()
	pdf.SetFillColor(240, 248, 255)
	pdf.Rect(30, currentY, 150, 70, "FD")

	pdf.SetXY(30, currentY+10)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "Period Summary", "", 1, "C", false, 0, "")

	pdf.SetXY(40, currentY+25)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Total Days:")
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, strconv.Itoa(len(reports)))

	pdf.SetXY(40, currentY+35)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Total Sales:")
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, strconv.Itoa(totalSalesCount))

	pdf.SetXY(40, currentY+45)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Total Revenue:")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 153, 51)
	pdf.Cell(0, 8, "Rs. "+strconv.FormatFloat(totalRevenue, 'f', 2, 64))

	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(40, currentY+55)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Total Discount:")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(204, 0, 0)
	pdf.Cell(0, 8, "Rs. "+strconv.FormatFloat(totalDiscount, 'f', 2, 64))

	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(50)

	// Daily breakdown table
	pdf.SetY(currentY + 85)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "Daily Breakdown", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(52, 73, 94)
	pdf.SetTextColor(255, 255, 255)

	colWidths := []float64{40, 30, 30, 40, 40}
	headers := []string{"Date", "Sales", "Discount", "Tax", "Revenue"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(8)

	// Table rows
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)

	for idx, report := range reports {
		if pdf.GetY() > 250 {
			pdf.AddPage()
			// Repeat header
			pdf.SetFont("Arial", "B", 9)
			pdf.SetFillColor(52, 73, 94)
			pdf.SetTextColor(255, 255, 255)
			for i, header := range headers {
				pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(8)
			pdf.SetFont("Arial", "", 8)
			pdf.SetTextColor(0, 0, 0)
		}

		if idx%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		rowData := []string{
			report.ReportDate.Format("Jan 2, 2006"),
			strconv.Itoa(report.TotalSales),
			"Rs. " + strconv.FormatFloat(report.TotalDiscount, 'f', 2, 64),
			"Rs. " + strconv.FormatFloat(report.TotalTax, 'f', 2, 64),
			"Rs. " + strconv.FormatFloat(report.TotalRevenue, 'f', 2, 64),
		}

		for i, data := range rowData {
			align := "L"
			if i > 0 {
				align = "R"
			}
			pdf.CellFormat(colWidths[i], 7, data, "1", 0, align, true, 0, "")
		}
		pdf.Ln(7)
	}

	// Footer
	pdf.SetY(-25)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 5, "System Generated Report", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")), "", 1, "C", false, 0, "")
}

func addDailyReportPage(pdf *gofpdf.Fpdf, report *dto.DailyReportDocument, pageNum, totalPages int) {
	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 12, fmt.Sprintf("Daily Report - %s", report.ReportDate.Format("January 2, 2006")), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 5, fmt.Sprintf("Page %d of %d", pageNum, totalPages), "", 1, "C", false, 0, "")
	pdf.Ln(8)

	// Sales overview box
	currentY := pdf.GetY()
	pdf.SetFillColor(245, 248, 250)
	pdf.Rect(15, currentY, 180, 45, "FD")
	pdf.SetXY(15, currentY+5)

	leftColX := 20.0
	rightColX := 105.0
	rowHeight := 12.0

	// Row 1
	currentY = pdf.GetY()
	pdf.SetXY(leftColX, currentY)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Total Sales:")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 102, 204)
	pdf.Cell(0, 6, strconv.Itoa(report.TotalSales))

	pdf.SetXY(rightColX, currentY)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(40, 6, "Total Revenue:")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 153, 51)
	pdf.Cell(0, 6, "Rs. "+strconv.FormatFloat(report.TotalRevenue, 'f', 2, 64))

	// Row 2
	currentY += rowHeight
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(leftColX, currentY)
	pdf.SetFont("Arial", "", 9)
	pdf.Cell(40, 6, "Cash Sales:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, strconv.Itoa(report.CashSales)+" (Rs. "+strconv.FormatFloat(report.CashRevenue, 'f', 2, 64)+")")

	pdf.SetXY(rightColX, currentY)
	pdf.SetFont("Arial", "", 9)
	pdf.Cell(40, 6, "Card Sales:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, strconv.Itoa(report.CardSales)+" (Rs. "+strconv.FormatFloat(report.CardRevenue, 'f', 2, 64)+")")

	// Row 3
	currentY += rowHeight
	pdf.SetXY(leftColX, currentY)
	pdf.SetFont("Arial", "", 9)
	pdf.Cell(40, 6, "Total Discount:")
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(204, 0, 0)
	pdf.Cell(0, 6, "Rs. "+strconv.FormatFloat(report.TotalDiscount, 'f', 2, 64))

	pdf.SetXY(rightColX, currentY)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(40, 6, "Total Tax:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, "Rs. "+strconv.FormatFloat(report.TotalTax, 'f', 2, 64))

	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(18)

	// Top selling items
	if len(report.TopSellingItems) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 8, "Top Selling Items", "", 1, "L", false, 0, "")
		pdf.Ln(2)

		// Table header
		pdf.SetFont("Arial", "B", 8)
		pdf.SetFillColor(52, 73, 94)
		pdf.SetTextColor(255, 255, 255)

		colWidths := []float64{15, 75, 25, 30, 35}
		headers := []string{"Rank", "Product", "Qty", "Price", "Total"}

		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 7, header, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(7)

		// Table rows
		pdf.SetFont("Arial", "", 7)
		pdf.SetTextColor(0, 0, 0)

		displayCount := len(report.TopSellingItems)
		if displayCount > 10 {
			displayCount = 10
		}

		for idx := 0; idx < displayCount; idx++ {
			item := report.TopSellingItems[idx]

			if idx%2 == 0 {
				pdf.SetFillColor(245, 245, 245)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}

			productName := item.ProductName
			if len(productName) > 45 {
				productName = productName[:42] + "..."
			}

			rowData := []string{
				strconv.Itoa(idx + 1),
				productName,
				strconv.Itoa(item.Quantity),
				"Rs. " + strconv.FormatFloat(item.UnitPrice, 'f', 2, 64),
				"Rs. " + strconv.FormatFloat(item.TotalAmount, 'f', 2, 64),
			}

			for i, data := range rowData {
				align := "L"
				if i == 0 {
					align = "C"
				}
				if i > 1 {
					align = "R"
				}
				pdf.CellFormat(colWidths[i], 6, data, "1", 0, align, true, 0, "")
			}
			pdf.Ln(6)
		}
	}

	// Footer
	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 7)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 4, fmt.Sprintf("Report Date: %s", report.ReportDate.Format("2006-01-02")), "", 1, "C", false, 0, "")
}
