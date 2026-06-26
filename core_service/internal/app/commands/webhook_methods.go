package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	kafkaMessages "github.com/testovoleg/5s-microservice-template/proto/kafka"
)

type WebhookMethodsCmdHandler interface {
	Webhook(ctx context.Context, command *WebhookCommand) error
}

type webhookMethodsHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
	redisRepo     repository.CacheRepository
	adminRepo     repository.AdminRepository
	storageRepo   repository.StorageRepository
	cloakRepo     repository.IDMRepository
}

func NewWebhookMethodsHandler(
	log logger.Logger,
	cfg *config.Config,
	kafkaProducer kafkaClient.Producer,
	redisRepo repository.CacheRepository,
	adminRepo repository.AdminRepository,
	storageRepo repository.StorageRepository,
	cloakRepo repository.IDMRepository,

) *webhookMethodsHandler {
	return &webhookMethodsHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
		redisRepo:     redisRepo,
		adminRepo:     adminRepo,
		storageRepo:   storageRepo,
		cloakRepo:     cloakRepo,
	}
}

type WebhookCommand struct {
	Payload *kafkaMessages.Payload
}
