package api

import (
	"bytes"
	"employee-crud/dao"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
)

// ExpiringStock represents a product/batch that is expiring within 3 months
// You may want to adjust this struct to match your actual DAO/DTO

// Use dao.ProductWithStockInfo directly

// Handler to generate PDF report of stocks expiring within 3 months
func GetExpiringStocksReportPDF(c *fiber.Ctx) error {
	// Use Sri Lanka timezone (UTC+5:30)
	sriLankaLoc := time.FixedZone("Asia/Colombo", 5*3600+30*60)
	now := time.Now().In(sriLankaLoc)
	threeMonthsLater := now.AddDate(0, 3, 0)

	// Fetch all products with stock (use DAO method with high limit, no cursor)
	const maxLimit = 10000
	stocks, _, _, err := dao.DB_FindAllProductsWithStock(maxLimit, "")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch stock data: " + err.Error(),
		})
	}

	// Filter stocks expiring within 3 months
	var expiringStocks []dao.ProductWithStockInfo
	for _, s := range stocks {
		if s.ExpiryDate == nil || s.ExpiryDate.IsZero() {
			continue
		}
		if s.ExpiryDate.After(now) && s.ExpiryDate.Before(threeMonthsLater) {
			expiringStocks = append(expiringStocks, s)
		}
	}

	// Generate PDF
	pdfBytes, err := generateExpiringStocksPDF(expiringStocks, now)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate PDF: " + err.Error(),
		})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=Expiring-Stocks-Report-%s.pdf", now.Format("2006-01-02")))
	c.Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	return c.Send(pdfBytes)
}

func generateExpiringStocksPDF(stocks []dao.ProductWithStockInfo, reportDate time.Time) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 15, "Expiring Stocks Report within the next 3 months", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(60, 60, 60)
	pdf.CellFormat(0, 8, "Report Date: "+reportDate.Format("Monday, January 2, 2006"), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	if len(stocks) == 0 {
		pdf.SetFont("Arial", "I", 12)
		pdf.SetTextColor(200, 0, 0)
		pdf.CellFormat(0, 10, "No stocks expiring within the next 3 months.", "", 1, "C", false, 0, "")
	} else {
		// Table header
		pdf.SetFont("Arial", "B", 10)
		pdf.SetFillColor(52, 73, 94)
		pdf.SetTextColor(255, 255, 255)
		colWidths := []float64{23, 50, 28, 18, 28, 23}
		headers := []string{"Product ID", "Product Name", "Batch ID", "Qty", "Expiry Date", "Status"}
		for i, h := range headers {
			pdf.CellFormat(colWidths[i], 8, h, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(8)

		// Table rows
		pdf.SetFont("Arial", "", 9)
		pdf.SetTextColor(0, 0, 0)
		for idx, s := range stocks {
			if pdf.GetY() > 260 {
				pdf.AddPage()
				pdf.SetFont("Arial", "B", 10)
				pdf.SetFillColor(52, 73, 94)
				pdf.SetTextColor(255, 255, 255)
				for i, h := range headers {
					pdf.CellFormat(colWidths[i], 8, h, "1", 0, "C", true, 0, "")
				}
				pdf.Ln(8)
				pdf.SetFont("Arial", "", 9)
				pdf.SetTextColor(0, 0, 0)
			}
			if idx%2 == 0 {
				pdf.SetFillColor(245, 245, 245)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			row := []string{
				s.ProductId,
				s.Name,
				s.BatchId,
				strconv.Itoa(s.StockQty),
				"",
				"",
			}
			if s.ExpiryDate != nil {
				row[4] = s.ExpiryDate.Format("2006-01-02")
			}
			row[5] = s.ProductStatus
			for i, data := range row {
				align := "L"
				if i == 3 { // Qty right align
					align = "R"
				} else if i == 4 || i == 5 { // Expiry Date and Status center align
					align = "C"
				}
				pdf.CellFormat(colWidths[i], 7, data, "1", 0, align, true, 0, "")
			}
			pdf.Ln(7)
		}
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 5, "This is a system generated document.", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Report generated on %s (UTC+5:30)", reportDate.Format("2006-01-02 15:04:05")), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}