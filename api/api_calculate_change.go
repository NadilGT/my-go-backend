package api

import (
	"github.com/gofiber/fiber/v2"
)

func CalculateChangeApi(c *fiber.Ctx) error {
	type CalculateChangeRequest struct {
		Total          float64 `json:"total" binding:"required"`
		AmountReceived float64 `json:"amountReceived" binding:"required"`
	}

	var req CalculateChangeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.AmountReceived < req.Total {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount received is less than total",
		})
	}

	change := req.AmountReceived - req.Total

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"change": change,
	})
}
