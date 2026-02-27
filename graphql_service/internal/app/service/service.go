package service

import (
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/app/commands"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type GraphQLService struct {
	Commands *commands.Commands
}

func NewGraphQLService(
	log logger.Logger,
	cfg *config.Config,
	kafkaProducer kafkaClient.Producer,
	coreClient coreService.CoreServiceClient,
) *GraphQLService {

	getTipicalDataHandler := commands.NewGetTipicalDataHandler(log, cfg, kafkaProducer, coreClient)
	postTipicalMutationHandler := commands.NewPostTipicalMutationHandler(log, cfg, kafkaProducer, coreClient)

	commands := commands.NewCommands(
		getTipicalDataHandler,
		postTipicalMutationHandler,
	)

	return &GraphQLService{Commands: commands}
}
