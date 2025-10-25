package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AddStockRequest represents the request body for adding stock
type AddStockRequest struct {
	ProductId    string     `json:"productId" validate:"required"`
	StockQty     int        `json:"stockQty" validate:"required,gt=0"`
	ExpiryDate   *time.Time `json:"expiryDate"`
	CostPrice    float64    `json:"costPrice" validate:"gt=0"`
	SellingPrice float64    `json:"sellingPrice" validate:"gt=0"`
}

// AddStock adds stock to an existing product
// If expiry date matches existing batch, adds to that batch
// If expiry date is different, creates new batch
func AddStock(c *fiber.Ctx) error {
	var req AddStockRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	// Add stock to product
	product, batchId, err := dao.DB_AddStockToProduct(
		req.ProductId,
		req.StockQty,
		req.ExpiryDate,
		req.CostPrice,
		req.SellingPrice,
	)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Sync stock to Stocks collection
	if err := dao.DB_SyncSingleProductStock(product); err != nil {
		// Log but don't fail
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Stock added successfully",
		"productId": product.ProductId,
		"batchId":   batchId,
		"product":   product,
	})
}
