package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type UpdateApiCmdHandler interface {
	Handle(ctx context.Context, command *UpdateApiCommand) (*dto.ApiDto, error)
}

type updateApiHandler struct {
	log        logger.Logger
	cfg        *config.Config
	coreClient coreService.CoreServiceClient
}

func NewUpdateApiHandler(log logger.Logger, cfg *config.Config, coreClient coreService.CoreServiceClient) *updateApiHandler {
	return &updateApiHandler{log: log, cfg: cfg, coreClient: coreClient}
}

func (c *updateApiHandler) Handle(ctx context.Context, command *UpdateApiCommand) (*dto.ApiDto, error) {
	ctx, span := tracing.StartSpan(ctx, "updateApiHandler.Handle")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	res, err := c.coreClient.UpdateApi(ctx, dto.UpdateApiToGrpcMessage(command.Params, command.Dto))
	if err != nil {
		return nil, err
	}

	return dto.ApiDtoFromGrpc(res), nil
}
