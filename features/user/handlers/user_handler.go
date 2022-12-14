package handlers

import (
	"github.com/bonkzero404/gaskn/database/stores"
	responseDto "github.com/bonkzero404/gaskn/dto"
	"github.com/bonkzero404/gaskn/features/user/dto"
	"github.com/bonkzero404/gaskn/features/user/interactors"
	"github.com/bonkzero404/gaskn/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService interactors.User
}

func NewUserHandler(userService interactors.User) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (service *UserHandler) CreateUser(c *fiber.Ctx) error {
	var request dto.UserCreateRequest
	var routeInternal = false

	if c.Params("CreateUser") == "create" {
		routeInternal = true
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ApiUnprocessableEntity(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:body-parser"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	errors := utils.ValidateStruct(request, c)
	if errors != nil {
		return utils.ApiErrorValidation(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:validate"),
			Cause:   utils.Lang(c, "global:err:validate-cause"),
			Inputs:  errors,
		})
	}

	response, err := service.UserService.CreateUser(c, &request, routeInternal)

	if err != nil {
		re := err.(*responseDto.ApiErrorResponse)
		return utils.ApiResponseError(c, re.StatusCode, responseDto.Errors{
			Message: utils.Lang(c, "user:err:create-failed"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	return utils.ApiCreated(c, response)
}

func (service *UserHandler) UserActivation(c *fiber.Ctx) error {
	var request dto.UserActivationRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.ApiUnprocessableEntity(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:body-parser"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	errors := utils.ValidateStruct(request, c)
	if errors != nil {
		return utils.ApiErrorValidation(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:validate"),
			Cause:   utils.Lang(c, "global:err:validate-cause"),
			Inputs:  errors,
		})
	}

	response, err := service.UserService.UserActivation(c, request.Code)

	if err != nil {
		re := err.(*responseDto.ApiErrorResponse)
		return utils.ApiResponseError(c, re.StatusCode, responseDto.Errors{
			Message: utils.Lang(c, "user:err:activation-failed"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	return utils.ApiCreated(c, response)
}

func (service *UserHandler) ReCreateUserActivation(c *fiber.Ctx) error {
	var request dto.UserReActivationRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.ApiUnprocessableEntity(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:body-parser"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	errors := utils.ValidateStruct(request, c)
	if errors != nil {
		return utils.ApiErrorValidation(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:validate"),
			Cause:   utils.Lang(c, "global:err:validate-cause"),
			Inputs:  errors,
		})
	}

	response, err := service.UserService.CreateUserActivation(c, request.Email, stores.ACTIVATION_CODE)

	if err != nil {
		re := err.(*responseDto.ApiErrorResponse)
		return utils.ApiResponseError(c, re.StatusCode, responseDto.Errors{
			Message: utils.Lang(c, "user:err:re-activation-failed"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	return utils.ApiCreated(c, response)
}

func (service *UserHandler) CreateActivationForgotPassword(c *fiber.Ctx) error {
	var request dto.UserForgotPassRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.ApiUnprocessableEntity(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:body-parser"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	errors := utils.ValidateStruct(request, c)
	if errors != nil {
		return utils.ApiErrorValidation(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:validate"),
			Cause:   utils.Lang(c, "global:err:validate-cause"),
			Inputs:  errors,
		})
	}

	response, err := service.UserService.CreateUserActivation(c, request.Email, stores.FORGOT_PASSWORD)

	if err != nil {
		re := err.(*responseDto.ApiErrorResponse)
		return utils.ApiResponseError(c, re.StatusCode, responseDto.Errors{
			Message: utils.Lang(c, "user:err:forgot-pass-failed"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	return utils.ApiCreated(c, response)
}

func (service *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	var request dto.UserForgotPassActRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.ApiUnprocessableEntity(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:body-parser"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	errors := utils.ValidateStruct(request, c)
	if errors != nil {
		return utils.ApiErrorValidation(c, responseDto.Errors{
			Message: utils.Lang(c, "global:err:validate"),
			Cause:   utils.Lang(c, "global:err:validate-cause"),
			Inputs:  errors,
		})
	}

	response, err := service.UserService.UpdatePassword(c, &request)

	if err != nil {
		re := err.(*responseDto.ApiErrorResponse)
		return utils.ApiResponseError(c, re.StatusCode, responseDto.Errors{
			Message: utils.Lang(c, "user:err:update-pass-failed"),
			Cause:   err.Error(),
			Inputs:  nil,
		})
	}

	return utils.ApiCreated(c, response)
}
