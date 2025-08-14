package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindAllEmployees(c *fiber.Ctx) error {
	returnValue, err := dao.DB_FindAllEmployees()
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(&returnValue)
}
