package commands

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type InvoiceHandlersListCmdHandler interface {
	Handle(ctx context.Context, command *InvoiceHandlersListCommand) ([]*dto.InvoiceHandlerDto, error)
}

type invoiceHandlersListHandler struct {
	log      logger.Logger
	cfg      *config.Config
	csClient coreService.CoreServiceClient
}

func NewInvoiceHandlersListHandler(log logger.Logger, cfg *config.Config, csClient coreService.CoreServiceClient) *invoiceHandlersListHandler {
	return &invoiceHandlersListHandler{log: log, cfg: cfg, csClient: csClient}
}

func (c *invoiceHandlersListHandler) Handle(ctx context.Context, command *InvoiceHandlersListCommand) ([]*dto.InvoiceHandlerDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "invoiceHandlersListHandler.Handle")
	defer span.Finish()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := c.csClient.InvoiceHandlersList(ctx, &coreService.InvoiceHandlersListReq{})
	if err != nil {
		return nil, err
	}

	return dto.InvoiceHandlerDtoListFromGrpc(res), nil
}
