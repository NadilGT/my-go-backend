package apiHandlers

import (
	"employee-crud/api"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Fiber on Render!")
	})

	app.Post("/employee", api.CreateUser)
	app.Get("/employee", api.FindAllEmployees)
	app.Put("/updateEmployee", api.UpdateEmployee)
	app.Get("/employee", api.FindEmployeeByID)
	app.Delete("/employee", api.SoftDeleteEmployeeById)
	app.Delete("/employee/hardDelete", api.HardDeleteEmployeeById)
}
