package role_assignment

import (
	"gaskn/database/driver"
	roleClientRepo "gaskn/features/role/repositories"
	"gaskn/features/role_assignment/handlers"
	"gaskn/features/role_assignment/repositories"
	"gaskn/features/role_assignment/services"

	"github.com/gofiber/fiber/v2"
)

// RegisterModule /*
func RegisterModule(app *fiber.App) {

	roleClientRepository := roleClientRepo.NewRoleClientRepository(driver.DB)
	roleAssignmentRepository := repositories.NewRoleAssignmentRepository(driver.Enforcer)
	roleAssignmentService := services.NewRoleAssignmentService(roleAssignmentRepository, roleClientRepository)
	RoleAssignmentHandler := handlers.NewRoleAssignmentHandler(roleAssignmentService)

	routesInit := ApiRouteClient{
		RoleAssignmentHandler: *RoleAssignmentHandler,
	}

	routesInit.Route(app)

}
