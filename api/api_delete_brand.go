package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteBrandApi(c *fiber.Ctx) error {
	id := c.Query("brandId")

	if id == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "Brand ID is required", nil)
	}

	if err := dao.DB_DeleteBrandByID(id); err != nil {
		return utils.NewCustomError(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Brand deleted successfully",
	})
}
