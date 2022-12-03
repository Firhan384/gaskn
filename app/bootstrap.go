package app

import (
	"gaskn/app/middleware"
	"gaskn/config"
	"gaskn/features/auth"
	cl "gaskn/features/client"
	"gaskn/features/role"
	"gaskn/features/role_assignment"
	"gaskn/features/user"
	"gaskn/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func features(app *fiber.App) {
	feature := utils.RouteFeature{}

	// Get feature lists
	app.Get(
		utils.SetupApiGroup()+"/features",
		middleware.Authenticate(),
		middleware.Permission(),
		func(c *fiber.Ctx) error {
			return utils.ApiOk(c, utils.ExtractRouteAsFeatures(c.App(), false))
		}).
		Name(
			feature.
				SetGroup("Features").
				SetName("FeatureLists").
				SetDescription("Admin get get route lists").
				SetOnlyAdmin(true).
				Exec(),
		)

	// Get feature per group
	app.Get(
		utils.SetupApiGroup()+"/features/group",
		middleware.Authenticate(),
		middleware.Permission(),
		func(c *fiber.Ctx) error {
			return utils.ApiOk(c, utils.FeaturesGroupLists(c.App(), false))
		}).
		Name(
			feature.
				SetGroup("Features").
				SetName("FeatureGroupLists").
				SetDescription("Admin get get group route lists").
				SetOnlyAdmin(true).
				Exec(),
		)

	if config.Config("TENANCY") == "true" {
		// Get feature lists Tenant
		app.Get(
			utils.SetupSubApiGroup()+"/features",
			middleware.Authenticate(),
			middleware.Permission(),
			func(c *fiber.Ctx) error {
				return utils.ApiOk(c, utils.ExtractRouteAsFeatures(c.App(), true))
			}).
			Name(
				feature.
					SetGroup("Client/Features").
					SetName("FeatureLists").
					SetDescription("Admin get get route lists").
					SetTenant(true).
					Exec(),
			)

		// Get feature per group tenant
		app.Get(
			utils.SetupSubApiGroup()+"/features/group",
			middleware.Authenticate(),
			middleware.Permission(),
			func(c *fiber.Ctx) error {
				return utils.ApiOk(c, utils.FeaturesGroupLists(c.App(), true))
			}).
			Name(
				feature.
					SetGroup("Client/Features").
					SetName("FeatureGroupLists").
					SetDescription("Admin get get group route lists").
					SetTenant(true).
					Exec(),
			)
	}
}

// Bootstrap /*
func Bootstrap(app *fiber.App) {
	// Monitor app
	app.Get("/monitor", monitor.New())

	// Register features
	features(app)

	// Register module user
	user.RegisterFeature(app)

	// Register module auth
	auth.RegisterFeature(app)

	// Register module role
	role.RegisterFeature(app)

	// Register Client
	cl.RegisterFeature(app)

	// Register feature Role Assignment
	role_assignment.RegisterFeature(app)
}

func SetupLogs() {
	if config.Config("ENABLE_LOG") == "true" {
		utils.CraeteDirectory(config.Config("LOG_LOCATION"))
	}
}
