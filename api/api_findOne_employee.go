package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindEmployeeByID(c *fiber.Ctx) error {
	//id := c.Params("id")
	id := c.Query("id")
	returnValue, err := dao.DB_FindEmployeeByID(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(&returnValue)
}
