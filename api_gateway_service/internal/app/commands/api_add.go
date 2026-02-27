package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type AddApiCmdHandler interface {
	Handle(ctx context.Context, command *AddApiCommand) (*dto.ApiDto, error)
}

type addApiHandler struct {
	log        logger.Logger
	cfg        *config.Config
	coreClient coreService.CoreServiceClient
}

func NewAddApiHandler(log logger.Logger, cfg *config.Config, coreClient coreService.CoreServiceClient) *addApiHandler {
	return &addApiHandler{log: log, cfg: cfg, coreClient: coreClient}
}

func (c *addApiHandler) Handle(ctx context.Context, command *AddApiCommand) (*dto.ApiDto, error) {
	ctx, span := tracing.StartSpan(ctx, "addApiHandler.Handle")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	res, err := c.coreClient.AddApi(ctx, dto.AddApiToGrpcMessage(command.Params, command.Dto))
	if err != nil {
		return nil, err
	}

	return dto.ApiDtoFromGrpc(res), nil
}
