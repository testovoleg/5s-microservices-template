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

type UpdateProductCmdHandler interface {
	Handle(ctx context.Context, command *UpdateProductCommand) error
}

type updateProductCmdHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewUpdateProductHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *updateProductCmdHandler {
	return &updateProductCmdHandler{log: log, cfg: cfg, kafkaProducer: kafkaProducer}
}

func (c *updateProductCmdHandler) Handle(ctx context.Context, command *UpdateProductCommand) error {
	ctx, span := tracing.StartSpan(ctx, "updateProductCmdHandler.Handle")
	defer span.End()

	updateDto := &kafkaMessages.ProductUpdate{
		ProductID:   command.UpdateDto.ProductID.String(),
		Name:        command.UpdateDto.Name,
		Description: command.UpdateDto.Description,
		Price:       command.UpdateDto.Price,
	}

	dtoBytes, err := proto.Marshal(updateDto)
	if err != nil {
		return err
	}

	return c.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductUpdate.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(ctx),
	})
}
