package app

import (
	"go-starterkit-project/config"
	"go-starterkit-project/modules/auth"
	"go-starterkit-project/modules/role"
	"go-starterkit-project/modules/user"
	"go-starterkit-project/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

/**
This function is used to register all modules,
this registration is the last process to register
all modules
*/
func Bootstrap(app *fiber.App) {
	// Monitor app
	app.Get("/monitor", monitor.New())

	// Register module user
	user.RegisterModule(app)

	// Register module auth
	auth.RegisterModule(app)

	// Register module role
	role.RegisterModule(app)
}

func SetupLogs() {
	if config.Config("ENABLE_LOG") == "true" {
		utils.CraeteDirectory(config.Config("LOG_LOCATION"))
	}
}
