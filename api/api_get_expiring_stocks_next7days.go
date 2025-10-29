package api

import (
	"employee-crud/dao"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler to get top 10 stocks expiring within next 7 days, sorted by highest stock quantity
func GetExpiringStocksNext7Days(c *fiber.Ctx) error {
	// Use Sri Lanka timezone (UTC+5:30)
	sriLankaLoc := time.FixedZone("Asia/Colombo", 5*3600+30*60)
	now := time.Now().In(sriLankaLoc)
	sevenDaysLater := now.AddDate(0, 0, 7)

	// Use optimized DAO method
	topN := 10
	expiringStocks, err := dao.DB_FindTopExpiringStocksNext7Days(topN, now, sevenDaysLater)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch expiring stocks: " + err.Error(),
		})
	}

	// Prepare response
	result := make([]fiber.Map, len(expiringStocks))
	for i, s := range expiringStocks {
		result[i] = fiber.Map{
			"product_id":     s.ProductId,
			"product_name":   s.Name,
			"batch_id":       s.BatchId,
			"stock_qty":      s.StockQty,
			"expiry_date":    "",
			"product_status": s.ProductStatus,
		}
		if s.ExpiryDate != nil {
			result[i]["expiry_date"] = s.ExpiryDate.Format("2006-01-02")
		}
	}

	return c.JSON(fiber.Map{
		"now":             now.Format("2006-01-02 15:04:05"),
		"timezone":        "+5:30 GMT",
		"expiring_stocks": result,
	})
}
