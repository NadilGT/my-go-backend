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

func GetMonthlyReturnsReportPDF(c *fiber.Ctx) error {
	monthStr := c.Query("month") // format: YYYY-MM
	if monthStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing month parameter (format: YYYY-MM)"})
	}
	// Parse month and set Sri Lanka timezone
	sriLankaLoc := time.FixedZone("Asia/Colombo", 5*3600+30*60)
	monthTime, err := time.ParseInLocation("2006-01", monthStr, sriLankaLoc)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid month format. Use YYYY-MM (e.g., 2026-08)"})
	}
	// Calculate start and end of month
	start := time.Date(monthTime.Year(), monthTime.Month(), 1, 0, 0, 0, 0, sriLankaLoc)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	returns, err := dao.GetReturnsByDateRange(c.Context(), start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch returns: " + err.Error()})
	}

	pdfBytes, err := generateMonthlyReturnsPDF(monthTime, returns)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate PDF: " + err.Error()})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=Returns-Report-%s.pdf", monthTime.Format("2006-01")))
	c.Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	return c.Send(pdfBytes)
}

// truncateString truncates a string to fit within specified width
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-2] + ".."
}

func generateMonthlyReturnsPDF(month time.Time, returns []dto.ReturnDTO) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)

	// Title
	pdf.SetFont("Arial", "B", 22)
	pdf.CellFormat(0, 15, "Monthly Returns Report", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 8, "Report Month: "+month.Format("January 2006"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, "Generated on: "+time.Now().Format("2006-01-02 15:04:05"), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Overview
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, fmt.Sprintf("Total Returns: %d", len(returns)), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Adjusted column widths to fit A4 page (total: 180mm usable width)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(52, 73, 94)
	pdf.SetTextColor(255, 255, 255)
	colWidths := []float64{22, 28, 24, 22, 20, 18, 26, 20}
	headers := []string{"Date", "Customer", "Contact", "Bill No.", "Product", "Amount", "Reason", "Notes"}
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)
	
	// Alternating row colors for better readability
	fillColor := false
	
	for _, ret := range returns {
		for _, prod := range ret.Products {
			if fillColor {
				pdf.SetFillColor(240, 240, 240)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			
			// Truncate long text to fit columns
			customerName := truncateString(ret.CustomerName, 18)
			billNumber := truncateString(ret.OriginalBillNumber, 14)
			productID := truncateString(prod.ProductID, 12)
			reason := truncateString(prod.Reason, 16)
			notes := truncateString(ret.AdditionalNotes, 12)
			
			pdf.CellFormat(colWidths[0], 7, ret.CreatedAt[:10], "1", 0, "C", fillColor, 0, "")
			pdf.CellFormat(colWidths[1], 7, customerName, "1", 0, "L", fillColor, 0, "")
			pdf.CellFormat(colWidths[2], 7, ret.ContactNumber, "1", 0, "L", fillColor, 0, "")
			pdf.CellFormat(colWidths[3], 7, billNumber, "1", 0, "L", fillColor, 0, "")
			pdf.CellFormat(colWidths[4], 7, productID, "1", 0, "L", fillColor, 0, "")
			pdf.CellFormat(colWidths[5], 7, fmt.Sprintf("%.2f", prod.Amount), "1", 0, "R", fillColor, 0, "")
			pdf.CellFormat(colWidths[6], 7, reason, "1", 0, "L", fillColor, 0, "")
			pdf.CellFormat(colWidths[7], 7, notes, "1", 0, "L", fillColor, 0, "")
			pdf.Ln(7)
			
			fillColor = !fillColor
		}
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 5, "This is a system generated document.", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Report covers returns from %s to %s (UTC+5:30)",
		month.Format("2006-01-02"), month.AddDate(0, 1, -1).Format("2006-01-02")), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}