package service

import (
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/app/commands"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type Service struct {
	Commands *commands.Commands
}

func NewAppService(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer, csClient coreService.CoreServiceClient) *Service {

	invoiceHandlersList := commands.NewInvoiceHandlersListHandler(log, cfg, csClient)
	updateProductHandler := commands.NewUpdateProductHandler(log, cfg, kafkaProducer)
	deleteProductHandler := commands.NewDeleteProductHandler(log, cfg, kafkaProducer)

	commands := commands.NewCommands(invoiceHandlersList, updateProductHandler, deleteProductHandler)

	return &Service{Commands: commands}
}
