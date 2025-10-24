package api

import (
	"context"
	"employee-crud/dao"
	"employee-crud/dto"
	"employee-crud/functions"
	"employee-crud/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProduct(c *fiber.Ctx) error {
	inputObj := dto.Product{}

	if err := c.BodyParser(&inputObj); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	ctx := context.Background()
	now := time.Now().UTC()

	// Check if a product with same category, brand, and subcategory exists
	existingProduct, err := dao.DB_FindProductByAttributes(
		inputObj.CategoryID,
		inputObj.BrandID,
		inputObj.SubCategoryID,
	)

	if err != nil && err != mongo.ErrNoDocuments {
		// Database error (not "not found" error)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if existingProduct != nil {
		// Product with same attributes exists
		// Check if expiry dates are different
		expiryDatesMatch := false
		if existingProduct.ExpiryDate == nil && inputObj.ExpiryDate == nil {
			expiryDatesMatch = true
		} else if existingProduct.ExpiryDate != nil && inputObj.ExpiryDate != nil {
			expiryDatesMatch = existingProduct.ExpiryDate.Equal(*inputObj.ExpiryDate)
		}

		if expiryDatesMatch {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "A product with the same category, brand, and subcategory already exists with the same expiry date")
		}

		// Different expiry date - create a new batch
		batchId, err := dao.GenerateId(ctx, "Batches", "BATCH")
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}

		newBatch := dto.Batch{
			BatchId:      batchId,
			StockQty:     inputObj.StockQty,
			ExpiryDate:   inputObj.ExpiryDate,
			CostPrice:    inputObj.CostPrice,
			SellingPrice: inputObj.SellingPrice,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// If existing product doesn't have batches array, initialize it with existing product data
		if len(existingProduct.Batches) == 0 {
			// Create first batch from existing product data
			firstBatchId, err := dao.GenerateId(ctx, "Batches", "BATCH")
			if err != nil {
				return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
			}

			firstBatch := dto.Batch{
				BatchId:      firstBatchId,
				StockQty:     existingProduct.StockQty,
				ExpiryDate:   existingProduct.ExpiryDate,
				CostPrice:    existingProduct.CostPrice,
				SellingPrice: existingProduct.SellingPrice,
				CreatedAt:    existingProduct.CreatedAt,
				UpdatedAt:    now,
			}

			// Update existing product to use batches
			existingProduct.Batches = []dto.Batch{firstBatch}
			if err := dao.DB_UpdateProductWithBatch(existingProduct, firstBatch); err != nil {
				return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
			}
		}

		// Add new batch to product
		if err := dao.DB_AddBatchToProduct(existingProduct.ProductId, newBatch); err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}

		// Sync stocks - fetch updated product first
		updatedProduct, err := dao.DB_FindProductById(existingProduct.ProductId)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}

		if err := dao.DB_SyncSingleProductStock(updatedProduct); err != nil {
			// Log but don't fail
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":   "Batch added to existing product successfully",
			"productId": existingProduct.ProductId,
			"batchId":   batchId,
		})
	}

	// No existing product found - create new product
	id, err := dao.GenerateId(ctx, "Products", "PRD")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	inputObj.ProductId = id
	inputObj.CreatedAt = now
	inputObj.UpdatedAt = now

	if err := functions.UniqueCheck(inputObj, "Products", []string{"ProductId"}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if validationErr := validate.Struct(inputObj); validationErr != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, validationErr.Error())
	}

	err = dao.DB_CreateProduct(&inputObj)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Automatically sync the product stock to Stocks collection
	if err := dao.DB_SyncSingleProductStock(&inputObj); err != nil {
		// Log the error but don't fail the product creation
		// You can add logging here if needed
		// For now, we'll silently continue
	}

	return utils.SendSuccessResponse(c)
}
