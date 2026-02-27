package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type GetFullApiCmdHandler interface {
	Handle(ctx context.Context, command *GetFullApiCommand) (*dto.ApiFullDto, error)
}

type getFullApiHandler struct {
	log        logger.Logger
	cfg        *config.Config
	coreClient coreService.CoreServiceClient
}

func NewGetFullApiHandler(log logger.Logger, cfg *config.Config, coreClient coreService.CoreServiceClient) *getFullApiHandler {
	return &getFullApiHandler{log: log, cfg: cfg, coreClient: coreClient}
}

func (c *getFullApiHandler) Handle(ctx context.Context, command *GetFullApiCommand) (*dto.ApiFullDto, error) {
	ctx, span := tracing.StartSpan(ctx, "getFullApiHandler.Handle")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	res, err := c.coreClient.GetFullApi(ctx, dto.GetFullApiToGrpcMessage(command.Params))
	if err != nil {
		return nil, err
	}

	return dto.ApiFullDtoFromGrpc(res), nil
}
