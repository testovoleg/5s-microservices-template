package service

import (
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/commands"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/queries"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type Service struct {
	Commands *commands.Commands
	Queries  *queries.ProductQueries
}

func NewAppService(
	log logger.Logger,
	cfg *config.Config,
	mongoRepo repository.Repository,
	redisRepo repository.CacheRepository,
) *Service {

	invoiceHandlersListCmdHandler := commands.NewInvoiceHandlersListHandler(log, cfg, mongoRepo, redisRepo)
	deleteProductCmdHandler := commands.NewDeleteProductCmdHandler(log, cfg, mongoRepo, redisRepo)
	updateProductCmdHandler := commands.NewUpdateProductCmdHandler(log, cfg, mongoRepo, redisRepo)

	getProductByIdHandler := queries.NewGetProductByIdHandler(log, cfg, mongoRepo, redisRepo)
	searchProductHandler := queries.NewSearchProductHandler(log, cfg, mongoRepo, redisRepo)

	commands := commands.NewCommands(invoiceHandlersListCmdHandler, updateProductCmdHandler, deleteProductCmdHandler)
	queries := queries.NewProductQueries(getProductByIdHandler, searchProductHandler)

	return &Service{Commands: commands, Queries: queries}
}
