package main

import (
	"employee-crud/apiHandlers"
	"employee-crud/dbConfigs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New()

	dbConfigs.ConnectMongoDB("mongodb+srv://admin:W6ptbj7HPS3RJ4cU@cluster0.tgypip5.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	apiHandlers.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))

}
