package main

import (
	"employee-crud/apiHandlers"
	"employee-crud/dbConfigs"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New()

	dbConfigs.ConnectMongoDB()

	apiHandlers.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + port))

}
