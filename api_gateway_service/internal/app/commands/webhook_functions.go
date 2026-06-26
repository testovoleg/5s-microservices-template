package commands

import (
	"context"
	"errors"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
	kafkaMessages "github.com/testovoleg/5s-microservice-template/proto/kafka"
	"google.golang.org/protobuf/proto"
)

func (c *webhookMethodsHandler) Webhook(ctx context.Context, command *WebhookCommand) error {
	ctx, span := tracing.StartSpan(ctx, "webhookMethodsHandler.Webhook")
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

func (c *webhookMethodsHandler) webhookToKafkaMessage(params *dto.ApiParamsDto, in []byte) (*kafkaMessages.Payload, error) {
	if in == nil {
		return nil, errors.New("invalid input")
	}

	return &kafkaMessages.Payload{
		EventName: constants.ContextKeyWebhook,
		CreatedAt: utils.CurrentTimestamppb(),
		Body:      &kafkaMessages.PayloadBody{Message: string(in)},
		Params:    dto.ApiParamsDtoToKafkaMessage(params),
	}, nil
}
