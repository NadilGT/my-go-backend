package api

import (
	"employee-crud/dao"
	"employee-crud/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UpdateGRNStatusRequest struct {
	GRNId  string `json:"grnId" validate:"required"`
	Status string `json:"status" validate:"required,oneof=pending completed partial_received"`
}

func UpdateGRNStatusApi(c *fiber.Ctx) error {
	var req UpdateGRNStatusRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	// Validate required fields
	if req.GRNId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "GRN ID is required")
	}

	if req.Status == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Status is required")
	}

	// Validate status values
	validStatuses := map[string]bool{
		"pending":          true,
		"completed":        true,
		"partial_received": true,
	}

	if !validStatuses[req.Status] {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid status. Must be one of: pending, completed, partial_received")
	}

	// Check if GRN exists
	exists, err := dao.DB_CheckGRNExists(req.GRNId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Error checking GRN existence")
	}

	if !exists {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "GRN not found")
	}

	// Update the status
	err = dao.DB_UpdateGRNStatus(req.GRNId, req.Status, time.Now().UTC())
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update GRN status")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "GRN status updated successfully",
		"grnId":   req.GRNId,
		"status":  req.Status,
	})
}
