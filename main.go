package main

import (
	"employee-crud/apiHandlers"
	"employee-crud/dbConfigs"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Configure CORS to allow your specific domains
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080,https://pos-frontend-tan.vercel.app",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,HEAD,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,Access-Control-Request-Method,Access-Control-Request-Headers",
		ExposeHeaders:    "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type",
		AllowCredentials: true,
	}))

	dbConfigs.ConnectMongoDB()

	// Setup TTL index for Sales collection (auto-delete after 24 hours)
	if err := dbConfigs.SetupSalesTTL(); err != nil {
		log.Fatal("Failed to setup Sales TTL index:", err)
	}

	apiHandlers.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + port))

}
