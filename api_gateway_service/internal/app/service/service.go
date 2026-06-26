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

	apiMethodsHandler := commands.NewApiMethodsHandler(log, cfg, coreClient)

	webhookMethodsHandler := commands.NewWebhookMethodsHandler(log, cfg, kafkaProducer)

	commands := commands.NewCommands(
		apiMethodsHandler,

		webhookMethodsHandler,
	)

	return &GatewayService{Commands: commands}
}
