package main

import (
	"employee-crud/apiHandlers"
	"employee-crud/dao"
	"employee-crud/dbConfigs"
	"employee-crud/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// GZIP Compression Middleware - Compresses responses for better performance
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Balance between speed and compression
	}))

	// Response Time Logger - Monitor API performance
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Configure CORS to allow your specific domains
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080,https://pos-frontend-tan.vercel.app",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,HEAD,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,Access-Control-Request-Method,Access-Control-Request-Headers",
		ExposeHeaders:    "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type",
		AllowCredentials: true,
	}))

	dbConfigs.ConnectMongoDB()
	dao.InitReturnsCollection(dbConfigs.DATABASE)

	// Setup TTL index for Sales collection (auto-delete after 24 hours)
	if err := dbConfigs.SetupSalesTTL(); err != nil {
		log.Fatal("Failed to setup Sales TTL index:", err)
	}

	// Setup TTL index for DailyReports collection (auto-delete at end of month)
	if err := dbConfigs.SetupDailyReportsTTL(); err != nil {
		log.Fatal("Failed to setup DailyReports TTL index:", err)
	}

	// Start background scheduler for automatic daily report saving
	utils.StartDailyReportScheduler()

	// Check and save any missing reports from the past 7 days
	go utils.SaveMissingReports()

	apiHandlers.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + port))

}
