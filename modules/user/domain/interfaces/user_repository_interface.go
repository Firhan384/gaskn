package interfaces

import (
	"go-boilerplate-clean-arch/domain/stores"
)

type UserRepositoryInterface interface {
	CreateUser(user *stores.User, userActivate *stores.UserActivation) (*stores.User, error)

	FindUserByEmail(email string) (*stores.User, error)

	FindUserById(id string) (*stores.User, error)

	FindUserActivationCode(userId string, code string) (*stores.UserActivation, error)

	UpdateUserActivation(id string, stat bool) (*stores.User, error)

	CreateUserActivation(userActivate *stores.UserActivation) (*stores.UserActivation, error)

	UpdatePassword(id string, password string) (*stores.User, error)

	UpdateActivationCodeUsed(userId string, code string) (*stores.UserActivation, error)
}
