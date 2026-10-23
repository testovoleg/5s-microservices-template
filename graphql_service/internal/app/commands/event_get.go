package commands

import (
	"context"

	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	model "github.com/testovoleg/5s-microservice-template/graphql_service/internal/graph_model"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type EventGetCommand struct {
	Params *model.GeneralParamsInput
	Input  *model.TipicalEventGetInput
}

type EventGetCmdHandler interface {
	Handle(ctx context.Context, command *EventGetCommand) (*model.TipicalEvent, error)
}

type eventGetHandler struct {
	log        logger.Logger
	cfg        *config.Config
	coreClient coreService.CoreServiceClient
}

func NewEventGetHandler(
	log logger.Logger,
	cfg *config.Config,
	coreClient coreService.CoreServiceClient,
) *eventGetHandler {
	return &eventGetHandler{
		log:        log,
		cfg:        cfg,
		coreClient: coreClient,
	}
}

func (c *eventGetHandler) Handle(ctx context.Context, command *EventGetCommand) (*model.TipicalEvent, error) {
	ctx, span := tracing.StartSpan(ctx, "eventGetHandler.Handle")
	defer span.End()

	grpcReq := model.EventGetInputToGrpc(command.Input)
	grpcReq.Params = model.GeneralParamsInputToGrpc(command.Params)

	res, err := c.coreClient.EventGet(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	return model.EventFromGrpc(res.GetData()), nil
}
