package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type ApiMethodsCmdHandler interface {
	AddApi(ctx context.Context, command *AddApiCommand) (*dto.ApiDto, error)
	GetApi(ctx context.Context, command *GetApiCommand) ([]*dto.ApiDto, error)
	GetFullApi(ctx context.Context, command *GetFullApiCommand) (*dto.ApiFullDto, error)
	UpdateApi(ctx context.Context, command *UpdateApiCommand) (*dto.ApiDto, error)
	DeleteApi(ctx context.Context, command *DeleteApiCommand) error
}

type apiMethodsHandler struct {
	log        logger.Logger
	cfg        *config.Config
	coreClient coreService.CoreServiceClient
}

func NewApiMethodsHandler(log logger.Logger, cfg *config.Config, coreClient coreService.CoreServiceClient) *apiMethodsHandler {
	return &apiMethodsHandler{log: log, cfg: cfg, coreClient: coreClient}
}

type AddApiCommand struct {
	Params *dto.ApiParamsDto
	Dto    *dto.AddApiReqDto
}

type GetApiCommand struct {
	Params *dto.ApiParamsDto
}

type GetFullApiCommand struct {
	Params *dto.ApiParamsDto
}

type UpdateApiCommand struct {
	Params *dto.ApiParamsDto
	Dto    *dto.UpdateApiReqDto
}

type DeleteApiCommand struct {
	Params *dto.ApiParamsDto
}
