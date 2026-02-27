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

	addApiHandler := commands.NewAddApiHandler(log, cfg, cloakRepo, adminRepo, redisRepo, kafkaProducer)
	getApiHandler := commands.NewGetApiHandler(log, cfg, cloakRepo, adminRepo, redisRepo)
	getFullApiHandler := commands.NewGetFullApiHandler(log, cfg, cloakRepo, adminRepo, redisRepo)
	updateApiHandler := commands.NewUpdateApiHandler(log, cfg, cloakRepo, adminRepo, redisRepo)
	deleteApiHandler := commands.NewDeleteApiHandler(log, cfg, cloakRepo, adminRepo, redisRepo)

	webhookHandler := commands.NewWebhookHandler(log, cfg, kafkaProducer, redisRepo, adminRepo, storageRepo, cloakRepo)

	commands := commands.NewCommands(
		addApiHandler,
		getApiHandler,
		getFullApiHandler,
		updateApiHandler,
		deleteApiHandler,

		webhookHandler,
	)

	return &CoreService{Commands: commands}
}
