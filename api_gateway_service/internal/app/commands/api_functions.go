package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

func (c *apiMethodsHandler) AddApi(ctx context.Context, command *AddApiCommand) (*dto.ApiDto, error) {
	ctx, span := tracing.StartSpan(ctx, "apiMethodsHandler.AddApi")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	res, err := c.coreClient.AddApi(ctx, dto.AddApiToGrpcMessage(command.Params, command.Dto))
	if err != nil {
		return nil, err
	}

	return dto.ApiDtoFromGrpc(res), nil
}

func (c *apiMethodsHandler) GetApi(ctx context.Context, command *GetApiCommand) ([]*dto.ApiDto, error) {
	ctx, span := tracing.StartSpan(ctx, "apiMethodsHandler.GetApi")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	res, err := c.coreClient.GetApi(ctx, dto.GetApiToGrpcMessage(command.Params))
	if err != nil {
		return nil, err
	}

	return dto.ListApiDtoFromGrpc(res.GetApi()), nil
}

func (c *apiMethodsHandler) GetFullApi(ctx context.Context, command *GetFullApiCommand) (*dto.ApiFullDto, error) {
	ctx, span := tracing.StartSpan(ctx, "apiMethodsHandler.GetFullApi")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	res, err := c.coreClient.GetFullApi(ctx, dto.GetFullApiToGrpcMessage(command.Params))
	if err != nil {
		return nil, err
	}

	return dto.ApiFullDtoFromGrpc(res), nil
}

func (c *apiMethodsHandler) UpdateApi(ctx context.Context, command *UpdateApiCommand) (*dto.ApiDto, error) {
	ctx, span := tracing.StartSpan(ctx, "apiMethodsHandler.UpdateApi")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	res, err := c.coreClient.UpdateApi(ctx, dto.UpdateApiToGrpcMessage(command.Params, command.Dto))
	if err != nil {
		return nil, err
	}

	return dto.ApiDtoFromGrpc(res), nil
}

func (c *apiMethodsHandler) DeleteApi(ctx context.Context, command *DeleteApiCommand) error {
	ctx, span := tracing.StartSpan(ctx, "apiMethodsHandler.DeleteApi")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	_, err := c.coreClient.DeleteApi(ctx, dto.DeleteApiToGrpcMessage(command.Params))
	if err != nil {
		return err
	}

	return nil
}
