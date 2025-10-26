package api

import (
	"employee-crud/dao"
	"employee-crud/dto"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateReturnApi(c *fiber.Ctx) error {
	var req dto.ReturnDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	if req.CustomerName == "" || req.ContactNumber == "" || len(req.Products) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "CustomerName, ContactNumber, and Products are required"})
	}
	id := uuid.New().String()
	req.ID = id
	req.CreatedAt = time.Now().Format(time.RFC3339)
	if err := dao.InsertReturn(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save return"})
	}
	return c.Status(fiber.StatusCreated).JSON(req)
}
