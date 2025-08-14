package api

import (
	"employee-crud/dao"
	"employee-crud/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateUser(c *fiber.Ctx) error {
	var employee dto.Employee
	if err := c.BodyParser(&employee); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	employee.ID = uuid.New().String()

	if err := dao.DB_CreateEmployee(employee); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create employee",
		})
	}
	return c.Status(200).JSON(employee)
}
