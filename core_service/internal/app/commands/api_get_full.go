package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type GetFullApiCmdHandler interface {
	Handle(ctx context.Context, command *GetFullApiCommand) (*models.ApiFull, error)
}

type getFullApiHandler struct {
	log       logger.Logger
	cfg       *config.Config
	cloakRepo repository.IDMRepository
	adminRepo repository.AdminRepository
	redisRepo repository.CacheRepository
}

func NewGetFullApiHandler(
	log logger.Logger,
	cfg *config.Config,
	cloakRepo repository.IDMRepository,
	adminRepo repository.AdminRepository,
	redisRepo repository.CacheRepository,
) *getFullApiHandler {
	return &getFullApiHandler{
		log:       log,
		cfg:       cfg,
		cloakRepo: cloakRepo,
		adminRepo: adminRepo,
		redisRepo: redisRepo,
	}
}

func (c *getFullApiHandler) Handle(ctx context.Context, command *GetFullApiCommand) (*models.ApiFull, error) {
	ctx, span := tracing.StartSpan(ctx, "getFullApiHandler.Handle")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, c.redisRepo, command.Params)
	if err != nil {
		return nil, err
	}

	api, err := getApiData(ctx, c.log, c.adminRepo, c.redisRepo, &models.ApiParams{AccessToken: command.Params.AccessToken, CompanyUuid: company.Uuid, ApiUuid: command.Params.ApiUuid})
	if err != nil {
		return nil, err
	}

	subscriptions := []string{}
	// subscriptions, err := c.targetRepo.GetSubscriptions(ctx, company, api)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "targetRepo.GetSubscriptions")
	// }

	return &models.ApiFull{
		Uuid:          api.Uuid,
		Title:         api.Title,
		Description:   api.Description,
		Token:         api.Token,
		Subscriptions: subscriptions,
	}, nil
}
