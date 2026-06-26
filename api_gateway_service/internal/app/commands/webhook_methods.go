package commands

import (
	"context"

	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"

	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type WebhookMethodsCmdHandler interface {
	Webhook(ctx context.Context, command *WebhookCommand) error
}

type webhookMethodsHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewWebhookMethodsHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *webhookMethodsHandler {
	return &webhookMethodsHandler{log: log, cfg: cfg, kafkaProducer: kafkaProducer}
}

type WebhookCommand struct {
	Params  *dto.ApiParamsDto
	Payload []byte
}
