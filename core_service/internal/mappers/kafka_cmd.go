package mappers

import (
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/commands"
	kafkaMessages "github.com/testovoleg/5s-microservice-template/proto/kafka"
)

func NewWebhookCommand(payload *kafkaMessages.Payload) *commands.WebhookCommand {
	return &commands.WebhookCommand{
		Payload: payload,
	}
}
