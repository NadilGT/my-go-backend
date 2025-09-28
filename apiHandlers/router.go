package apiHandlers

import (
	"employee-crud/api"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Fiber on Render!")
	})

	app.Post("/CreateCategory", api.CreateCategoryApi)
	app.Get("/api/categories", api.FindAllCategoriesApi)

}
