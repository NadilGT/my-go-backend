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

	dbConfigs.ConnectMongoDB("mongodb+srv://admin:W6ptbj7HPS3RJ4cU@cluster0.tgypip5.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	apiHandlers.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + port))

}
