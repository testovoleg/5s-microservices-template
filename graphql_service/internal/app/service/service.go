package service

import (
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/app/mutations"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/app/queries"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type BugsService struct {
	Mutations *mutations.Mutations
	Queries   *queries.Queries
}

func NewBugsService(
	log logger.Logger,
	cfg *config.Config,
	kafkaProducer kafkaClient.Producer,
	coreClient coreService.CoreServiceClient,
) *BugsService {

	createBugHandler := mutations.NewCreateBugHandler(log, cfg, kafkaProducer, coreClient)

	getBugsHandler := queries.NewGetBugsHandler(log, cfg, coreClient)

	m := mutations.NewMutations(createBugHandler)
	q := queries.NewQueries(getBugsHandler)

	return &BugsService{Mutations: m, Queries: q}
}
