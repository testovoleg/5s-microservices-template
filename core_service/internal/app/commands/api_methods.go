package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type ApiMethodsCmdHandler interface {
	AddApi(ctx context.Context, command *AddApiCommand) (*models.Api, error)
	GetApi(ctx context.Context, command *GetApiCommand) ([]*models.Api, error)
	GetFullApi(ctx context.Context, command *GetFullApiCommand) (*models.ApiFull, error)
	UpdateApi(ctx context.Context, command *UpdateApiCommand) (*models.Api, error)
	DeleteApi(ctx context.Context, command *DeleteApiCommand) error
}

type apiMethodsHandler struct {
	log           logger.Logger
	cfg           *config.Config
	cloakRepo     repository.IDMRepository
	adminRepo     repository.AdminRepository
	redisRepo     repository.CacheRepository
	kafkaProducer kafkaClient.Producer
}

func NewApiMethodsHandler(
	log logger.Logger,
	cfg *config.Config,
	cloakRepo repository.IDMRepository,
	adminRepo repository.AdminRepository,
	redisRepo repository.CacheRepository,
	kafkaProducer kafkaClient.Producer,
) *apiMethodsHandler {
	return &apiMethodsHandler{
		log:           log,
		cfg:           cfg,
		cloakRepo:     cloakRepo,
		adminRepo:     adminRepo,
		redisRepo:     redisRepo,
		kafkaProducer: kafkaProducer,
	}
}

type AddApiCommand struct {
	Params *models.ApiParams
	Data   *models.Api
}

type GetApiCommand struct {
	Params *models.ApiParams
}

type GetFullApiCommand struct {
	Params *models.ApiParams
}

type UpdateApiCommand struct {
	Params *models.ApiParams
	Data   *models.UpdateApi
}

type DeleteApiCommand struct {
	Params *models.ApiParams
}
