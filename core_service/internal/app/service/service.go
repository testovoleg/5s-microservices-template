package service

import (
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/commands"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type Service struct {
	Commands *commands.Commands
}

func NewAppService(
	log logger.Logger,
	cfg *config.Config,
	redisRepo repository.CacheRepository,
) *Service {

	invoiceHandlersListCmdHandler := commands.NewInvoiceHandlersListHandler(log, cfg, redisRepo)
	deleteProductCmdHandler := commands.NewDeleteProductCmdHandler(log, cfg, redisRepo)
	updateProductCmdHandler := commands.NewUpdateProductCmdHandler(log, cfg, redisRepo)

	commands := commands.NewCommands(invoiceHandlersListCmdHandler, updateProductCmdHandler, deleteProductCmdHandler)

	return &Service{Commands: commands}
}
