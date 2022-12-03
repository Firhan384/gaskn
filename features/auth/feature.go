package auth

import (
	"gaskn/driver"
	"gaskn/features/auth/handlers"
	"gaskn/features/auth/services"
	roleRepo "gaskn/features/role/repositories"
	"gaskn/features/user/repositories"

	"github.com/gofiber/fiber/v2"
)

// RegisterFeature /*
func RegisterFeature(app *fiber.App) {
	userRepository := repositories.NewUserRepository(driver.DB)
	roleRepository := roleRepo.NewRoleRepository(driver.DB)
	authService := services.NewAuthService(userRepository, roleRepository)
	authHandler := handlers.NewAuthHandler(authService)

	routesInit := ApiRoute{
		AuthHandler: *authHandler,
	}

	routesInit.Route(app)
}
