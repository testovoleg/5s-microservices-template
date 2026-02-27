package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type DeleteApiCmdHandler interface {
	Handle(ctx context.Context, command *DeleteApiCommand) error
}

type deleteApiHandler struct {
	log        logger.Logger
	cfg        *config.Config
	coreClient coreService.CoreServiceClient
}

func NewDeleteApiHandler(log logger.Logger, cfg *config.Config, coreClient coreService.CoreServiceClient) *deleteApiHandler {
	return &deleteApiHandler{log: log, cfg: cfg, coreClient: coreClient}
}

func (c *deleteApiHandler) Handle(ctx context.Context, command *DeleteApiCommand) error {
	ctx, span := tracing.StartSpan(ctx, "deleteApiHandler.Handle")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	_, err := c.coreClient.DeleteApi(ctx, dto.DeleteApiToGrpcMessage(command.Params))
	if err != nil {
		return err
	}

	return nil
}
