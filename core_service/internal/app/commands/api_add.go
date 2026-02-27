package commands

import (
	"context"

	"github.com/pkg/errors"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

type AddApiCmdHandler interface {
	Handle(ctx context.Context, command *AddApiCommand) (*models.Api, error)
}

type addApiHandler struct {
	log           logger.Logger
	cfg           *config.Config
	cloakRepo     repository.IDMRepository
	adminRepo     repository.AdminRepository
	redisRepo     repository.CacheRepository
	kafkaProducer kafkaClient.Producer
}

func NewAddApiHandler(
	log logger.Logger,
	cfg *config.Config,
	cloakRepo repository.IDMRepository,
	adminRepo repository.AdminRepository,
	redisRepo repository.CacheRepository,
	kafkaProducer kafkaClient.Producer,
) *addApiHandler {
	return &addApiHandler{
		log:           log,
		cfg:           cfg,
		cloakRepo:     cloakRepo,
		adminRepo:     adminRepo,
		redisRepo:     redisRepo,
		kafkaProducer: kafkaProducer,
	}
}

func (c *addApiHandler) Handle(ctx context.Context, command *AddApiCommand) (*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "addApiHandler.Handle")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, c.redisRepo, command.Params)
	if err != nil {
		return nil, err
	}

	listApi, err := c.adminRepo.GetProperty(ctx, command.Params.AccessToken, company.Uuid)
	if err != nil {
		return nil, errors.Wrap(err, "adminRepo.GetProperty")
	}

	c.redisRepo.PutApiList(ctx, company.Uuid, listApi)

	api, err := c.createApi(command)
	if err != nil {
		return nil, errors.Wrap(err, "handle.createApi")
	}

	listApi = append(listApi, api)

	err = c.adminRepo.AddProperty(ctx, command.Params.AccessToken, company.Uuid, listApi)
	if err != nil {
		return nil, errors.Wrap(err, "adminRepo.AddProperty")
	}

	c.redisRepo.PutApiList(ctx, company.Uuid, listApi)

	return api, nil
}

func (c *addApiHandler) createApi(command *AddApiCommand) (*models.Api, error) {
	if command.Data == nil {
		return nil, errors.New("something went wrong")
	}

	api := &models.Api{
		Uuid:        utils.GenerateUuid(),
		Title:       command.Data.Title,
		Description: command.Data.Description,
		Token:       command.Data.Token,
	}

	if api.Token == "" || api.Token == "string" {
		return nil, errors.New("token is not filled")
	}

	if api.Title == "" || api.Title == "string" {
		api.Title = "Template Api"
	}

	if api.Description == "" || api.Description == "string" {
		api.Description = "Template Api"
	}

	return api, nil
}
