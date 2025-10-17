package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func GetCompletedGRNsCount(c *fiber.Ctx) error {
	count, err := dao.DB_CountGRNsByStatus("completed")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"completed_grns": count,
	})
}

func GetPendingGRNsCount(c *fiber.Ctx) error {
	count, err := dao.DB_CountGRNsByStatus("pending")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"pending_grns": count,
	})
}

func GetPartialReceivedGRNsCount(c *fiber.Ctx) error {
	count, err := dao.DB_CountGRNsByStatus("partial_received")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"partial_received_grns": count,
	})
}
