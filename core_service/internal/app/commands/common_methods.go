package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type EventGetCommand struct {
	Params *models.ApiParams
	Data   *models.EventGetReq
}

type CommonMethodsCmdHandler interface {
	EventGet(ctx context.Context, cmd *EventGetCommand) (*models.Event, error)
}

type commonMethodsHandler struct {
	log       logger.Logger
	cfg       *config.Config
	cloakRepo repository.IDMRepository
	adminRepo repository.AdminRepository
	redisRepo repository.CacheRepository
}

func NewCommonMethodsHandler(
	log logger.Logger,
	cfg *config.Config,
	cloakRepo repository.IDMRepository,
	adminRepo repository.AdminRepository,
	redisRepo repository.CacheRepository,
) *commonMethodsHandler {
	return &commonMethodsHandler{
		log:       log,
		cfg:       cfg,
		cloakRepo: cloakRepo,
		adminRepo: adminRepo,
		redisRepo: redisRepo,
	}
}

// EventGet reads the current status of an async event from Redis.
// The caller is authenticated via getUserData like any other method.
func (c *commonMethodsHandler) EventGet(ctx context.Context, cmd *EventGetCommand) (*models.Event, error) {
	ctx, span := tracing.StartSpan(ctx, "commonMethodsHandler.EventGet")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, c.redisRepo, cmd.Params)
	if err != nil {
		return nil, attrErr(span, cmd.Params, cmd.Data, err)
	}

	res, err := c.redisRepo.GetEvent(ctx, company.Uuid, cmd.Data.EventUuid)
	if err != nil {
		return nil, attrErr(span, cmd.Params, cmd.Data, err)
	}
	return res, nil
}
