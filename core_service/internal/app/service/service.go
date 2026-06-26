package service

import (
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/commands"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type CoreService struct {
	Commands *commands.Commands
}

func NewCoreService(
	log logger.Logger,
	cfg *config.Config,
	cloakRepo repository.IDMRepository,
	adminRepo repository.AdminRepository,
	redisRepo repository.CacheRepository,
	kafkaProducer kafkaClient.Producer,
	storageRepo repository.StorageRepository,
) *CoreService {

	apiMethodsHandler := commands.NewApiMethodsHandler(log, cfg, cloakRepo, adminRepo, redisRepo, kafkaProducer)

	webhookMethodsHandler := commands.NewWebhookMethodsHandler(log, cfg, kafkaProducer, redisRepo, adminRepo, storageRepo, cloakRepo)

	commands := commands.NewCommands(
		apiMethodsHandler,

		webhookMethodsHandler,
	)

	return &CoreService{Commands: commands}
}
