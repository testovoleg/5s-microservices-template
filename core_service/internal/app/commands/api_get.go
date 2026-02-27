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

type GetApiCmdHandler interface {
	Handle(ctx context.Context, command *GetApiCommand) ([]*models.Api, error)
}

type getApiHandler struct {
	log       logger.Logger
	cfg       *config.Config
	cloakRepo repository.IDMRepository
	adminRepo repository.AdminRepository
	redisRepo repository.CacheRepository
}

func NewGetApiHandler(
	log logger.Logger,
	cfg *config.Config,
	cloakRepo repository.IDMRepository,
	adminRepo repository.AdminRepository,
	redisRepo repository.CacheRepository,
) *getApiHandler {
	return &getApiHandler{
		log:       log,
		cfg:       cfg,
		cloakRepo: cloakRepo,
		adminRepo: adminRepo,
		redisRepo: redisRepo,
	}
}

func (c *getApiHandler) Handle(ctx context.Context, command *GetApiCommand) ([]*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "getApiHandler.Handle")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, c.redisRepo, command.Params)
	if err != nil {
		return nil, err
	}

	listApi, err := c.redisRepo.GetApiList(ctx, company.Uuid)
	if err != nil || len(listApi) == 0 {
		listApi, err = c.adminRepo.GetProperty(ctx, command.Params.AccessToken, company.Uuid)
		if err != nil {
			return nil, errors.Wrap(err, "adminRepo.GetProperty")
		}
		c.redisRepo.PutApiList(ctx, company.Uuid, listApi)
	}

	return listApi, nil
}
