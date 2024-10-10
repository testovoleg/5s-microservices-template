package commands

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	kafkaMessages "github.com/testovoleg/5s-microservice-template/proto/kafka"
	"google.golang.org/protobuf/proto"
)

type DeleteProductCmdHandler interface {
	Handle(ctx context.Context, command *DeleteProductCommand) error
}

type deleteProductHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewDeleteProductHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *deleteProductHandler {
	return &deleteProductHandler{log: log, cfg: cfg, kafkaProducer: kafkaProducer}
}

func (c *deleteProductHandler) Handle(ctx context.Context, command *DeleteProductCommand) error {
	ctx, span := tracing.StartSpan(ctx, "deleteProductHandler.Handle")
	defer span.End()

	createDto := &kafkaMessages.ProductDelete{ProductID: command.ProductID.String()}

	dtoBytes, err := proto.Marshal(createDto)
	if err != nil {
		return err
	}

	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductDelete.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(ctx),
	})
}
