package interactors

import (
	"github.com/bonkzero404/gaskn/features/auth/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type UserAuth interface {
	Authenticate(c *fiber.Ctx, auth *dto.UserAuthRequest) (*dto.UserAuthResponse, error)

	GetProfile(c *fiber.Ctx, id string) (*dto.UserAuthProfileResponse, error)

	RefreshToken(c *fiber.Ctx, token *jwt.Token) (*dto.UserAuthResponse, error)
}
