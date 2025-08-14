package api

import (
	"employee-crud/dao"
	"employee-crud/dto"

	"github.com/gofiber/fiber/v2"
)

func UpdateEmployee(c *fiber.Ctx) error {
	var employee dto.Employee
	if err := c.BodyParser(&employee); err != nil {
		return err
	}
	if err := dao.DB_UpdateEmployee(&employee); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update employee",
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(employee)
}
