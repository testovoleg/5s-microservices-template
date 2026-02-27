package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type GetApiCmdHandler interface {
	Handle(ctx context.Context, command *GetApiCommand) ([]*dto.ApiDto, error)
}

type getApiHandler struct {
	log        logger.Logger
	cfg        *config.Config
	coreClient coreService.CoreServiceClient
}

func NewGetApiHandler(log logger.Logger, cfg *config.Config, coreClient coreService.CoreServiceClient) *getApiHandler {
	return &getApiHandler{log: log, cfg: cfg, coreClient: coreClient}
}

func (c *getApiHandler) Handle(ctx context.Context, command *GetApiCommand) ([]*dto.ApiDto, error) {
	ctx, span := tracing.StartSpan(ctx, "getApiHandler.Handle")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	res, err := c.coreClient.GetApi(ctx, dto.GetApiToGrpcMessage(command.Params))
	if err != nil {
		return nil, err
	}

	return dto.ListApiDtoFromGrpc(res.GetApi()), nil
}
