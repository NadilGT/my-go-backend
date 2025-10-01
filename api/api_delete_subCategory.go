package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteSubCategoryApi(c *fiber.Ctx) error {
	id := c.Query("subCategoryId")

	if id == "" {
		return utils.NewCustomError(c, fiber.StatusBadRequest, "SubCategory ID is required", nil)
	}

	if err := dao.DB_DeleteSubCategoryByID(id); err != nil {
		return utils.NewCustomError(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SubCategory deleted successfully",
	})
}
