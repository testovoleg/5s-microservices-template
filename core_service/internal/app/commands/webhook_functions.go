package commands

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
	kafkaMessages "github.com/testovoleg/5s-microservice-template/proto/kafka"
)

func (c *webhookMethodsHandler) Webhook(ctx context.Context, command *WebhookCommand) error {
	ctx, span := tracing.StartSpan(ctx, "webhookMethodsHandler.Webhook")
	defer span.End()

	if !c.isValidCommand(command) {
		return errors.New("invalid input")
	}

	params := command.Payload.Params

	company := &models.Company{Uuid: params.CompanyUuid}
	utils.Attr(span, "company_uuid", company.Uuid)

	adminToken, err := getAdminToken(ctx, c.cloakRepo, c.redisRepo)
	if err != nil {
		return err
	}

	api, err := getApiData(ctx, c.adminRepo, c.redisRepo, &models.ApiParams{
		AccessToken: adminToken, CompanyUuid: company.Uuid, ApiUuid: params.ApiUuid,
	})
	if err != nil {
		return err
	}

	webhookMessage := command.Payload.Body.Message
	if webhookMessage == "" {
		return errors.New("empty message")
	}

	var msg models.WebhookMessage
	err = json.Unmarshal([]byte(webhookMessage), &msg)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal update")
	}

	event, err := c.createKafkaEvent(ctx, api, company.Uuid, adminToken, msg)
	if err != nil {
		return errors.Wrap(err, "createKafkaEvent")
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	_ = eventBytes //delete this line
	// c.kafkaProducer.PublishMessageInTopic(ctx, c.cfg.KafkaTopics.ExportChatWebhook.TopicName, eventBytes)

	return nil
}

func (c *webhookMethodsHandler) isValidCommand(command *WebhookCommand) bool {
	return command != nil && command.Payload != nil && command.Payload.Params != nil && command.Payload.Body != nil
}

func (c *webhookMethodsHandler) createKafkaEvent(ctx context.Context, api *models.Api, companyUuid, adminToken string, msg models.WebhookMessage) (*kafkaMessages.Event, error) {
	out := &kafkaMessages.Event{
		EventUUID: utils.GenerateUuid(),
		Type:      kafkaMessages.EventType_MESSAGE,
		Gateway: &kafkaMessages.Gateway{
			CompanyUUID: companyUuid,
			Channel:     kafkaMessages.Channel_EMPTY_CHANNEL.Enum(),
			BotUUID:     api.Uuid,
			ChatID:      "",
		},
		Content: &kafkaMessages.Content{
			Body: &kafkaMessages.EventBody{},
		},
		CreatedAt: time.Now().Format(time.RFC3339),
		From:      &kafkaMessages.Contact{},
	}

	switch msg.Type {
	case "incomingMessage":
		// err := c.processContact(ctx, out, api, companyUuid, adminToken, msg)
		// if err != nil {
		// 	return nil, err
		// }

		// c.processIncomingMsg(ctx, out, companyUuid, adminToken, msg)

		if msg.Type == "outgoingMessage" {
			out.Direction = kafkaMessages.Direction_OUTCOMING
			out.From.IsBot = true
		}
		return out, nil

	case "messageStatus":
		// c.processStatuses(ctx, out, msg)
		return out, nil

	default:
		return nil, errors.New("unknown type: " + msg.Type)
	}
}
