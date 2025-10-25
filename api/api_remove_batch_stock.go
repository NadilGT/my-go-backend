package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

// RemoveStockRequest represents the request to remove stock from a batch
type RemoveStockRequest struct {
	ProductId        string `json:"productId" validate:"required"`
	BatchId          string `json:"batchId" validate:"required"`
	QuantityToRemove int    `json:"quantityToRemove" validate:"required,gt=0"`
}

// DeleteBatchRequest represents the request to delete a batch
type DeleteBatchRequest struct {
	ProductId string `json:"productId" validate:"required"`
	BatchId   string `json:"batchId" validate:"required"`
}

// RemoveStockFromBatch reduces stock from a specific batch
func RemoveStockFromBatch(c *fiber.Ctx) error {
	var req RemoveStockRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	product, err := dao.DB_RemoveStockFromBatch(req.ProductId, req.BatchId, req.QuantityToRemove)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Sync stock to Stocks collection
	if err := dao.DB_SyncSingleProductStock(product); err != nil {
		// Log but don't fail
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Stock removed successfully",
		"product": product,
	})
}

// DeleteBatch completely removes a batch from a product
func DeleteBatch(c *fiber.Ctx) error {
	var req DeleteBatchRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	product, err := dao.DB_DeleteBatch(req.ProductId, req.BatchId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Sync stock to Stocks collection
	if err := dao.DB_SyncSingleProductStock(product); err != nil {
		// Log but don't fail
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Batch deleted successfully",
		"product": product,
	})
}
