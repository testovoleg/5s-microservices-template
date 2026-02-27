package service

import (
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/app/commands"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type GatewayService struct {
	Commands *commands.Commands
}

func NewGatewayService(
	log logger.Logger, cfg *config.Config,
	kafkaProducer kafkaClient.Producer,
	coreClient coreService.CoreServiceClient,
) *GatewayService {

	addApiHandler := commands.NewAddApiHandler(log, cfg, coreClient)
	getApiHandler := commands.NewGetApiHandler(log, cfg, coreClient)
	getFullApiHandler := commands.NewGetFullApiHandler(log, cfg, coreClient)
	updateApiHandler := commands.NewUpdateApiHandler(log, cfg, coreClient)
	deleteApiHandler := commands.NewDeleteApiHandler(log, cfg, coreClient)

	webhookHandler := commands.NewWebhookHandler(log, cfg, kafkaProducer)

	commands := commands.NewCommands(
		addApiHandler,
		getApiHandler,
		getFullApiHandler,
		updateApiHandler,
		deleteApiHandler,

		webhookHandler,
	)

	return &GatewayService{Commands: commands}
}
