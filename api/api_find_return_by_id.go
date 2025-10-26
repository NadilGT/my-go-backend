package api

import (
	"employee-crud/dao"

	"github.com/gofiber/fiber/v2"
)

func FindReturnByIdApi(c *fiber.Ctx) error {
	id := c.Params("id")
	ret, err := dao.GetReturnByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Return not found"})
	}
	return c.JSON(ret)
}
