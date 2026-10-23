package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type EventGetCommand struct {
	Params *dto.ApiParamsDto
	Dto    *dto.EventGetReqDto
}

type CommonMethodsCmdHandler interface {
	EventGet(ctx context.Context, command *EventGetCommand) (*dto.EventDto, error)
}

type commonMethodsHandler struct {
	log        logger.Logger
	cfg        *config.Config
	coreClient coreService.CoreServiceClient
}

func NewCommonMethodsHandler(log logger.Logger, cfg *config.Config, coreClient coreService.CoreServiceClient) *commonMethodsHandler {
	return &commonMethodsHandler{log: log, cfg: cfg, coreClient: coreClient}
}

func (c *commonMethodsHandler) EventGet(ctx context.Context, command *EventGetCommand) (*dto.EventDto, error) {
	ctx, span := tracing.StartSpan(ctx, "commonMethodsHandler.EventGet")
	defer span.End()

	res, err := c.coreClient.EventGet(ctx, dto.EventGetReqDtoToGrpcMessage(command.Params, command.Dto))
	if err != nil {
		return nil, err
	}

	return dto.EventDtoFromGrpc(res.GetData()), nil
}
