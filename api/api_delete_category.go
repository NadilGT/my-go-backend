package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteCategoryApi(c *fiber.Ctx) error {
	id := c.Query("categoryId")

	if id == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "Category ID is required", nil)
	}

	if err := dao.DB_DeletecategoryByID(id); err != nil {
		return utils.NewCustomError(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
