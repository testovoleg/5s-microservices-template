package commands

import (
	"context"

	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	model "github.com/testovoleg/5s-microservice-template/graphql_service/internal/graph_model"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

type GetTipicalDataCmdHandler interface {
	Handle(ctx context.Context, command *GetTipicalDataCommand) (*model.TipicalQueryResponse, error)
}

type getTipicalDataHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
	coreClient    coreService.CoreServiceClient
}

func NewGetTipicalDataHandler(
	log logger.Logger,
	cfg *config.Config,
	kafkaProducer kafkaClient.Producer,
	coreClient coreService.CoreServiceClient,
) *getTipicalDataHandler {
	return &getTipicalDataHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
		coreClient:    coreClient,
	}
}

func (c *getTipicalDataHandler) Handle(ctx context.Context, command *GetTipicalDataCommand) (*model.TipicalQueryResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "getTipicalDataHandler.Handle")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	// res, err := c.coreClient.GetApi(ctx, model.GetApiToGrpcMessage(command.Params, command.Dto))
	// if err != nil {
	// 	return nil, err
	// }

	// return model.ApiFromGrpc(res), nil
	return &model.TipicalQueryResponse{UUID: utils.StrPtr(utils.GenerateUuid())}, nil
}
