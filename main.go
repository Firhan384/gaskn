package main

import (
	"fmt"
	appRoute "gaskn/app"
	"gaskn/config"
	"gaskn/database"
	"gaskn/database/driver"
	"gaskn/seeders"
	"gaskn/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
)

func main() {
	// Fiber app
	app := fiber.New()

	utils.SetupLang()

	// Setup Logs
	appRoute.SetupLogs()

	// Call database connection
	driver.ConnectDB()

	// Auto migration table
	database.MigrateDB()

	// Init Casbin
	driver.InitCasbin()

	// Handling global cors
	app.Use(cors.New())

	// Securing with helmet
	app.Use(helmet.New())

	// Call bootstrap all module
	appRoute.Bootstrap(app)

	// Run Seeder
	for _, seed := range seeders.All() {
		if err := seed.Run(driver.DB); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}

	// Set port
	appPort := fmt.Sprintf("%s:%s", config.Config("APP_HOST"), config.Config("APP_PORT"))

	// Listen app
	app.Listen(appPort)
}
