package commands

import (
	"context"

	"github.com/pkg/errors"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type UpdateApiCmdHandler interface {
	Handle(ctx context.Context, command *UpdateApiCommand) (*models.Api, error)
}

type updateApiHandler struct {
	log       logger.Logger
	cfg       *config.Config
	cloakRepo repository.IDMRepository
	adminRepo repository.AdminRepository
	redisRepo repository.CacheRepository
}

func NewUpdateApiHandler(
	log logger.Logger,
	cfg *config.Config,
	cloakRepo repository.IDMRepository,
	adminRepo repository.AdminRepository,
	redisRepo repository.CacheRepository,
) *updateApiHandler {
	return &updateApiHandler{
		log:       log,
		cfg:       cfg,
		cloakRepo: cloakRepo,
		adminRepo: adminRepo,
		redisRepo: redisRepo,
	}
}

func (c *updateApiHandler) Handle(ctx context.Context, command *UpdateApiCommand) (*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "updateApiHandler.Handle")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, c.redisRepo, command.Params)
	if err != nil {
		return nil, err
	}

	if command.Params.ApiUuid == "" {
		return nil, errors.New("apiUuid is not filled")
	}

	//do not use redis just in case
	listApi, err := c.adminRepo.GetProperty(ctx, command.Params.AccessToken, company.Uuid)
	if err != nil {
		return nil, errors.Wrap(err, "adminRepo.GetProperty")
	}

	api, updated := c.updateApi(listApi, command.Data, command.Params.ApiUuid)
	if !updated {
		return nil, errors.New("apiUuid not found")
	}

	err = c.adminRepo.AddProperty(ctx, command.Params.AccessToken, company.Uuid, listApi)
	if err != nil {
		return nil, errors.Wrap(err, "adminRepo.AddProperty")
	}

	c.redisRepo.PutApiList(ctx, company.Uuid, listApi)

	return api, nil
}

func (c *updateApiHandler) updateApi(list []*models.Api, update *models.UpdateApi, apiUuid string) (api *models.Api, updated bool) {
	setValue := func(exist, data *string) {
		if data != nil && *data != "string" {
			exist = data
		}
	}

	for _, a := range list {
		if a != nil && a.Uuid == apiUuid {
			setValue(&a.Title, update.Title)
			setValue(&a.Description, update.Description)
			setValue(&a.Token, update.Token)
			api = a
			updated = true
		}
	}

	return
}
