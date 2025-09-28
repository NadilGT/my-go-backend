package apiHandlers

import (
	"employee-crud/api"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Fiber on Render!")
	})

	app.Post("/api/createCategory", api.CreateCategoryApi)
	app.Get("/api/categories", api.FindAllCategoriesApi)
	app.Post("/api/brands", api.CreateBrand)
	app.Get("/api/brands", api.FindAllBrands)

}
