package commands

import (
	"context"
	"errors"
	"time"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	kafkaMessages "github.com/testovoleg/5s-microservice-template/proto/kafka"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WebhookCmdHandler interface {
	Handle(ctx context.Context, command *WebhookCommand) error
}

type webhookHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewWebhookHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *webhookHandler {
	return &webhookHandler{log: log, cfg: cfg, kafkaProducer: kafkaProducer}
}

func (c *webhookHandler) Handle(ctx context.Context, command *WebhookCommand) error {
	ctx, span := tracing.StartSpan(ctx, "WebhookHandler.Handle")
	defer span.End()

	msg, err := c.webhookToKafkaMessage(command.Params, command.Payload)
	if err != nil {
		return err
	}

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	err = c.kafkaProducer.PublishMessageInTopic(ctx, c.cfg.KafkaTopics.WebhookExample.TopicName, msgBytes)
	if err != nil {
		return err
	}

	return nil
}

func (c *webhookHandler) webhookToKafkaMessage(params *dto.ApiParamsDto, in []byte) (*kafkaMessages.Payload, error) {
	if in == nil {
		return nil, errors.New("Messsage body is empty")
	}

	return &kafkaMessages.Payload{
		EventName: "webhook",
		CreatedAt: timestamppb.New(time.Now()),
		Body: &kafkaMessages.PayloadBody{
			Message: string(in),
		},
		Params: dto.ApiParamsDtoToKafkaMessage(params),
	}, nil
}
