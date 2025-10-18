package api

import (
	"employee-crud/dao"
	"employee-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func FindSaleByIdApi(c *fiber.Ctx) error {
	saleId := c.Query("saleId")
	if saleId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "saleId parameter is required")
	}

	sale, err := dao.DB_FindSaleById(saleId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Sale not found")
	}

	return c.JSON(fiber.Map{
		"sale": sale,
	})
}

func FindAllSalesApi(c *fiber.Ctx) error {
	sales, err := dao.DB_FindAllSales()
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(sales)
}
