package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func GetBrandCostSummaryApi(c *fiber.Ctx) error {
	brandId := c.Query("brandId")

	if brandId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "brandId is required")
	}

	totalCost, expectedCost, err := dao.DB_GetBrandCostSummary(brandId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"brand_id":      brandId,
		"total_spend":   totalCost,
		"sales_target":  expectedCost,
		"target_profit": expectedCost - totalCost,
	})
}
