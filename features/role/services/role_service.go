package services

import (
	"gaskn/database/stores"
	respModel "gaskn/dto"
	"gaskn/features/role/contracts"
	"gaskn/features/role/dto"
	"gaskn/utils"

	"github.com/gofiber/fiber/v2"
)

type RoleService struct {
	RoleRepository contracts.RoleRepository
}

func NewRoleService(
	roleRepository contracts.RoleRepository,
) contracts.RoleService {
	return &RoleService{
		RoleRepository: roleRepository,
	}
}

func (service RoleService) CreateRole(c *fiber.Ctx, role *dto.RoleRequest) (*dto.RoleResponse, error) {

	roleData := stores.Role{
		RoleName:        role.RoleName,
		RoleDescription: role.RoleDescription,
		IsActive:        true,
	}

	err := service.RoleRepository.CreateRole(&roleData).Error

	if err != nil {
		return &dto.RoleResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "global:err:failed-unknown"),
		}
	}

	roleResponse := dto.RoleResponse{
		ID:              roleData.ID.String(),
		RoleName:        roleData.RoleName,
		RoleDescription: roleData.RoleDescription,
		IsActive:        roleData.IsActive,
	}

	return &roleResponse, nil
}

func (service RoleService) GetRoleList(c *fiber.Ctx) (*utils.Pagination, error) {
	var roles []stores.Role
	var resp []dto.RoleResponse

	res, err := service.RoleRepository.GetRoleList(&roles, c)

	if err != nil {
		return nil, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    err.Error(),
		}
	}

	for _, item := range roles {
		resp = append(resp, dto.RoleResponse{
			ID:              item.ID.String(),
			RoleName:        item.RoleName,
			RoleDescription: item.RoleDescription,
			IsActive:        item.IsActive,
		})
	}

	res.Rows = resp

	return res, nil
}

func (service RoleService) UpdateRole(c *fiber.Ctx, id string, role *dto.RoleRequest) (*dto.RoleResponse, error) {
	// Check role if exists
	var roleStore stores.Role

	errCheckRole := service.RoleRepository.GetRoleById(&roleStore, id).Error

	if errCheckRole != nil {
		return &dto.RoleResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "role:err:read-exists"),
		}
	}

	roleStore.RoleName = role.RoleName
	roleStore.RoleDescription = role.RoleDescription
	roleStore.IsActive = true

	err := service.RoleRepository.UpdateRoleById(&roleStore).Error

	if err != nil {
		return &dto.RoleResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "global:err:failed-unknown"),
		}
	}

	roleResponse := dto.RoleResponse{
		ID:              roleStore.ID.String(),
		RoleName:        roleStore.RoleName,
		RoleDescription: roleStore.RoleDescription,
		IsActive:        roleStore.IsActive,
	}

	return &roleResponse, nil
}

func (service RoleService) DeleteRoleById(c *fiber.Ctx, id string) (*dto.RoleResponse, error) {
	// Check role if exists
	var roleStore stores.Role

	errCheckRole := service.RoleRepository.GetRoleById(&roleStore, id).Error

	if errCheckRole != nil {
		return &dto.RoleResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "role:err:read-exists"),
		}
	}

	err := service.RoleRepository.DeleteRoleById(&roleStore).Error

	if err != nil {
		return &dto.RoleResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "global:err:failed-unknown"),
		}
	}

	roleResponse := dto.RoleResponse{
		ID:              roleStore.ID.String(),
		RoleName:        roleStore.RoleName,
		RoleDescription: roleStore.RoleDescription,
		IsActive:        roleStore.IsActive,
	}

	return &roleResponse, nil
}