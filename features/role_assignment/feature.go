package role_assignment

import (
	"gaskn/driver"
	roleClientRepo "gaskn/features/role/repositories"
	"gaskn/features/role_assignment/handlers"
	"gaskn/features/role_assignment/services"

	"github.com/gofiber/fiber/v2"
)

// RegisterFeature /*
func RegisterFeature(app *fiber.App) {
	roleRepository := roleClientRepo.NewRoleRepository(driver.DB)
	roleClientRepository := roleClientRepo.NewRoleClientRepository(driver.DB)

	roleAssignmentService := services.NewRoleAssignmentService(roleClientRepository, roleRepository)
	RoleAssignmentHandler := handlers.NewRoleAssignmentHandler(roleAssignmentService)

	routesInit := ApiRouteClient{
		RoleAssignmentHandler: *RoleAssignmentHandler,
	}

	routesInit.Route(app)

}
