package commands

import (
	"context"

	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type PostTipicalMutationCmdHandler interface {
	Handle(ctx context.Context, command *PostTipicalMutationCommand) error
}

type postTipicalMutationHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
	coreClient    coreService.CoreServiceClient
}

func NewPostTipicalMutationHandler(
	log logger.Logger,
	cfg *config.Config,
	kafkaProducer kafkaClient.Producer,
	coreClient coreService.CoreServiceClient,
) *postTipicalMutationHandler {
	return &postTipicalMutationHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
		coreClient:    coreClient,
	}
}

func (c *postTipicalMutationHandler) Handle(ctx context.Context, command *PostTipicalMutationCommand) error {
	ctx, span := tracing.StartSpan(ctx, "postTipicalMutationHandler.Handle")
	defer span.End()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.SpanContext())
	// err := c.coreClient.DeleteApi(ctx, model.DeleteApiToGrpcMessage(command.Params))
	// if err != nil {
	// 	return err
	// }

	return nil
}
