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

type DeleteApiCmdHandler interface {
	Handle(ctx context.Context, command *DeleteApiCommand) error
}

type deleteApiHandler struct {
	log       logger.Logger
	cfg       *config.Config
	cloakRepo repository.IDMRepository
	adminRepo repository.AdminRepository
	redisRepo repository.CacheRepository
}

func NewDeleteApiHandler(
	log logger.Logger,
	cfg *config.Config,
	cloakRepo repository.IDMRepository,
	adminRepo repository.AdminRepository,
	redisRepo repository.CacheRepository,

) *deleteApiHandler {
	return &deleteApiHandler{
		log:       log,
		cfg:       cfg,
		cloakRepo: cloakRepo,
		adminRepo: adminRepo,
		redisRepo: redisRepo,
	}
}

func (c *deleteApiHandler) Handle(ctx context.Context, command *DeleteApiCommand) error {
	ctx, span := tracing.StartSpan(ctx, "deleteApiHandler.Handle")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, c.redisRepo, command.Params)
	if err != nil {
		return err
	}

	if command.Params.ApiUuid == "" {
		return errors.New("apiUuid is not filled")
	}

	listApi, err := c.adminRepo.GetProperty(ctx, command.Params.AccessToken, company.Uuid)
	if err != nil {
		return errors.Wrap(err, "adminRepo.GetProperty")
	}

	res, deleted := c.filterApiList(listApi, command.Params.ApiUuid)
	if !deleted {
		return errors.New("apiUuid not found")
	}

	//delete before add if admin change with error
	c.redisRepo.DeleteApiList(ctx, company)

	err = c.adminRepo.AddProperty(ctx, command.Params.AccessToken, company.Uuid, res)
	if err != nil {
		return errors.Wrap(err, "adminRepo.AddProperty")
	}

	return nil
}

func (c *deleteApiHandler) filterApiList(list []*models.Api, apiUuid string) (res []*models.Api, deleted bool) {
	for _, api := range list {
		if api.Uuid == apiUuid {
			deleted = true
			continue
		}

		res = append(res, api)
	}
	return
}
