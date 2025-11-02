package api

import (
	"context"
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

type UpdateSupplierStatusRequest struct {
	SupplierId string `json:"supplierId"`
	Status     string `json:"status"`
}

func UpdateSupplierStatus(c *fiber.Ctx) error {
	var req UpdateSupplierStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Status != "active" && req.Status != "inactive" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status value"})
	}
	if req.SupplierId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "SupplierId is required"})
	}
	if err := dao.DB_UpdateSupplierStatus(context.Background(), req.SupplierId, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Supplier status updated successfully"})
}
