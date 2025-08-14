package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func SoftDeleteEmployeeById(c *fiber.Ctx) error {
	//id := c.Params("id")
	id := c.Query("id")
	if err := dao.DB_SoftDeleteEmployeeByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete employee",
		})
	}
	return c.Status(fiber.StatusNoContent).SendString("Employee soft deleted successfully")
}

func HardDeleteEmployeeById(c *fiber.Ctx) error {
	//id := c.Params("id")
	id := c.Query("id")
	if err := dao.DB_HardDeleteEmployeeByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete employee",
		})
	}
	return c.Status(fiber.StatusNoContent).SendString("Employee hard deleted successfully")
}
