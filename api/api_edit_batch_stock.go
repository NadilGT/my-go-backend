package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// EditBatchStockRequest represents the request to edit batch stock quantity
type EditBatchStockRequest struct {
	ProductId string `json:"productId" validate:"required"`
	BatchId   string `json:"batchId" validate:"required"`
	StockQty  int    `json:"stockQty" validate:"gte=0"`
}

// EditBatchDetailsRequest represents the request to edit batch details
type EditBatchDetailsRequest struct {
	ProductId    string     `json:"productId" validate:"required"`
	BatchId      string     `json:"batchId" validate:"required"`
	ExpiryDate   *time.Time `json:"expiryDate"`
	CostPrice    float64    `json:"costPrice"`
	SellingPrice float64    `json:"sellingPrice"`
}

// EditBatchStock updates the stock quantity of a specific batch
func EditBatchStock(c *fiber.Ctx) error {
	var req EditBatchStockRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	product, err := dao.DB_EditBatchStock(req.ProductId, req.BatchId, req.StockQty)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Sync stock to Stocks collection
	if err := dao.DB_SyncSingleProductStock(product); err != nil {
		// Log but don't fail
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Batch stock updated successfully",
		"product": product,
	})
}

// EditBatchDetails updates batch details like expiry date and prices
func EditBatchDetails(c *fiber.Ctx) error {
	var req EditBatchDetailsRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	product, err := dao.DB_EditBatchDetails(
		req.ProductId,
		req.BatchId,
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
		"message": "Batch details updated successfully",
		"product": product,
	})
}
