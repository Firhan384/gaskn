package services

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gosimple/slug"

	"gaskn/config"
	"gaskn/database/driver"
	"gaskn/database/stores"
	respModel "gaskn/dto"
	"gaskn/features/client/contracts"
	"gaskn/features/client/dto"
	userContract "gaskn/features/user/contracts"
	"gaskn/utils"
)

type ClientService struct {
	ClientRepository contracts.ClientRepository
	UserRepository   userContract.UserRepository
}

func NewClientService(
	clientRepository contracts.ClientRepository,
	UserRepository userContract.UserRepository,
) contracts.ClientService {
	return &ClientService{
		ClientRepository: clientRepository,
		UserRepository:   UserRepository,
	}
}

func (service ClientService) CreateClient(c *fiber.Ctx, client *dto.ClientRequest, userId string) (*dto.ClientResponse, error) {
	pUuid, _ := uuid.Parse(userId)
	clientRoute := config.Config("API_WRAP") + "/" + config.Config("API_VERSION") + "/" + config.Config("API_CLIENT") + "/:" + config.Config("API_CLIENT_PARAM")

	clientStore := stores.Client{
		ClientName:        client.ClientName,
		ClientDescription: client.ClientDescription,
		ClientSlug:        slug.Make(client.ClientName),
		UserId:            pUuid,
		IsActive:          true,
	}

	// Create client
	role, err := service.ClientRepository.CreateClient(&clientStore)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return &dto.ClientResponse{}, &respModel.ApiErrorResponse{
				StatusCode: fiber.StatusUnprocessableEntity,
				Message:    utils.Lang(c, "client:err:duplicate"),
			}
		}

		return &dto.ClientResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, err.Error()),
		}
	}

	// Get User By Id
	var user = stores.User{}
	service.UserRepository.FindUserById(&user, pUuid.String())

	// Crete group policy
	if g, err := driver.AddGroupingPolicy(
		pUuid.String(),
		role.ID.String(),
		clientStore.ID.String(),
		user.FullName,
		role.RoleName,
		clientStore.ClientName,
	); !g {
		return &dto.ClientResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, err.Error()),
		}
	}

	// Create permission user to group
	if p, err := driver.AddPolicy(
		role.ID.String(),
		clientStore.ID.String(),
		"/"+clientRoute+"/*",
		"GET|POST|PUT|DELETE",
		"",
		role.RoleName,
		clientStore.ClientName,
	); !p {
		return &dto.ClientResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, err.Error()),
		}
	}

	roleResponse := dto.ClientResponse{
		ID:                clientStore.ID.String(),
		ClientName:        clientStore.ClientName,
		ClientDescription: clientStore.ClientDescription,
		ClientSlug:        clientStore.ClientSlug,
		IsActive:          clientStore.IsActive,
	}

	return &roleResponse, nil
}

func (service ClientService) UpdateClient(c *fiber.Ctx, client *dto.ClientRequest) (*dto.ClientResponse, error) {
	// Check role if exists
	var clientStore stores.Client

	clientId := c.Params(config.Config("API_CLIENT_PARAM"))

	errCheckClient := service.ClientRepository.GetClientById(&clientStore, clientId).Error

	if errCheckClient != nil {
		return &dto.ClientResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "client:err:read-exists"),
		}
	}

	clientStore.ClientName = client.ClientName
	clientStore.ClientDescription = client.ClientDescription
	clientStore.ClientSlug = slug.Make(client.ClientName)
	clientStore.IsActive = true

	err := service.ClientRepository.UpdateClientById(&clientStore).Error

	if err != nil {
		return &dto.ClientResponse{}, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    utils.Lang(c, "global:err:failed-unknown"),
		}
	}

	clientResponse := dto.ClientResponse{
		ID:                clientStore.ID.String(),
		ClientName:        clientStore.ClientName,
		ClientDescription: clientStore.ClientDescription,
		ClientSlug:        clientStore.ClientSlug,
		IsActive:          clientStore.IsActive,
	}

	return &clientResponse, nil
}

func (service ClientService) GetClientByUser(c *fiber.Ctx, userId string) (*utils.Pagination, error) {
	var clientAssignment []stores.ClientAssignment
	var resp []dto.ClientResponse

	res, err := service.ClientRepository.GetClientListByUser(&clientAssignment, c, userId)

	if err != nil {
		return nil, &respModel.ApiErrorResponse{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    err.Error(),
		}
	}

	for _, item := range clientAssignment {
		resp = append(resp, dto.ClientResponse{
			ID:                item.Client.ID.String(),
			ClientName:        item.Client.ClientName,
			ClientDescription: item.Client.ClientDescription,
			ClientSlug:        item.Client.ClientSlug,
			IsActive:          item.IsActive,
		})
	}

	res.Rows = resp

	return res, nil

}