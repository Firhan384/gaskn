package services

import (
	"errors"
	respModel "go-starterkit-project/domain/dto"
	"go-starterkit-project/domain/stores"
	"go-starterkit-project/modules/user/domain/dto"
	"go-starterkit-project/modules/user/domain/interfaces"
	"go-starterkit-project/modules/user/services/factories"
	"go-starterkit-project/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository           interfaces.UserRepositoryInterface
	UserActivationRepository interfaces.UserActivationRepositoryInterface
	RepositoryAggregate      interfaces.RepositoryAggregateInterface
	ActionFactory            factories.ActionFactoryInterface
}

func NewUserService(
	userRepository interfaces.UserRepositoryInterface,
	userActivationRepository interfaces.UserActivationRepositoryInterface,
	repositoryAggregate interfaces.RepositoryAggregateInterface,
	factory factories.ActionFactoryInterface,
) interfaces.UserServiceInterface {
	return &UserService{
		UserRepository:           userRepository,
		UserActivationRepository: userActivationRepository,
		RepositoryAggregate:      repositoryAggregate,
		ActionFactory:            factory,
	}
}

func (service UserService) CreateUser(c *fiber.Ctx, user *dto.UserCreateRequest) (*dto.UserCreateResponse, error) {
	hashPassword, _ := utils.HashPassword(user.Password)

	userData := stores.User{
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: hashPassword,
	}

	activationCode := utils.StringWithCharset(32)

	userAvtivate := stores.UserActivation{
		Code:    activationCode,
		ActType: stores.ACTIVATION_CODE,
	}

	result, err := service.RepositoryAggregate.CreateUser(&userData, &userAvtivate)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return &dto.UserCreateResponse{}, &respModel.ApiErrorResponse{
				StatusCode: fiber.StatusUnprocessableEntity,
				Message:    utils.Lang(c, "user:err:create:register:failed"),
			}
		}

		return &dto.UserCreateResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "global:err:failed-unknown"),
		}
	}

	sendMail := respModel.Mail{
		To:           []string{user.Email},
		Subject:      "User Activation",
		TemplateHtml: "user_activation.html",
		BodyParam: map[string]interface{}{
			"Name": user.FullName,
			"Code": activationCode,
		},
	}

	utils.SendMail(&sendMail)

	response := dto.UserCreateResponse{
		ID:       userData.ID.String(),
		FullName: result.FullName,
		Email:    result.Email,
		Phone:    result.Phone,
		IsActive: userData.IsActive,
	}

	return &response, nil
}

func (service UserService) UserActivation(c *fiber.Ctx, email string, code string) (*dto.UserCreateResponse, error) {
	var user stores.User
	var userAct stores.UserActivation

	errUser := service.UserRepository.FindUserByEmail(&user, email).Error

	if errors.Is(errUser, gorm.ErrRecordNotFound) {
		return &dto.UserCreateResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusNotFound,
			Message:    utils.Lang(c, "user:err:create:user-not-found"),
		}
	}

	if user.IsActive {
		return &dto.UserCreateResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "user:err:activate:already-active"),
		}
	}

	errAct := service.UserActivationRepository.FindUserActivationCode(&userAct, user.ID.String(), code).Error

	if errors.Is(errAct, gorm.ErrRecordNotFound) {
		return &dto.UserCreateResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusNotFound,
			Message:    utils.Lang(c, "user:err:create:activation-not-found"),
		}
	}

	t := time.Now()

	if userAct.ExpiredAt.Before(t) {
		return &dto.UserCreateResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusGone,
			Message:    utils.Lang(c, "user:err:create:activation-expired"),
		}
	}

	userNew, errUserNew := service.RepositoryAggregate.UpdateUserActivation(user.ID.String(), true)

	if errUserNew != nil {
		return &dto.UserCreateResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    errUserNew.Error(),
		}
	}

	service.RepositoryAggregate.UpdateActivationCodeUsed(user.ID.String(), code)

	response := dto.UserCreateResponse{
		ID:       userNew.ID.String(),
		FullName: userNew.FullName,
		Email:    userNew.Email,
		Phone:    userNew.Phone,
		IsActive: userNew.IsActive,
	}

	return &response, nil
}

func (service UserService) CreateUserActivation(c *fiber.Ctx, email string, actType stores.ActivationType) (map[string]interface{}, error) {
	var user stores.User

	errUser := service.UserRepository.FindUserByEmail(&user, email).Error

	if errors.Is(errUser, gorm.ErrRecordNotFound) {
		return nil, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusNotFound,
			Message:    utils.Lang(c, "user:err:create:user-not-found"),
		}
	}

	if user.IsActive && actType == stores.ACTIVATION_CODE {
		return nil, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "user:err:activate:already-active"),
		}
	}

	_, errActFactory := service.ActionFactory.Create(actType, &user)

	if errActFactory != nil {
		return nil, errActFactory
	}

	return map[string]interface{}{}, nil
}

func (service UserService) UpdatePassword(c *fiber.Ctx, forgotPassReq *dto.UserForgotPassActRequest) (map[string]interface{}, error) {
	var user stores.User
	var userAct stores.UserActivation

	if forgotPassReq.Password != forgotPassReq.RepeatPassword {
		return nil, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "user:err:pass-match"),
		}
	}

	errUser := service.UserRepository.FindUserByEmail(&user, forgotPassReq.Email).Error

	if errors.Is(errUser, gorm.ErrRecordNotFound) {
		return nil, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusNotFound,
			Message:    utils.Lang(c, "user:err:create:user-not-found"),
		}
	}

	errAct := service.UserActivationRepository.FindUserActivationCode(&userAct, user.ID.String(), forgotPassReq.Code).Error

	if errors.Is(errAct, gorm.ErrRecordNotFound) {
		return nil, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusNotFound,
			Message:    utils.Lang(c, "user:err:create:activation-not-found"),
		}
	}

	if userAct.IsUsed {
		return nil, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "user:err:pass-code-used"),
		}
	}

	t := time.Now()

	if userAct.ExpiredAt.Before(t) {
		return nil, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusGone,
			Message:    utils.Lang(c, "user:err:create:activation-expired"),
		}
	}

	go func() {
		hashPassword, _ := utils.HashPassword(user.Password)

		userData := stores.User{
			FullName: user.FullName,
			Email:    user.Email,
			Phone:    user.Phone,
			Password: hashPassword,
		}

		service.UserRepository.UpdatePassword(&userData)
		service.RepositoryAggregate.UpdateActivationCodeUsed(user.ID.String(), forgotPassReq.Code)
	}()

	return map[string]interface{}{}, nil
}
