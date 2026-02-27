package commands

import "github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"

type WebhookCommand struct {
	Params  *dto.ApiParamsDto
	Payload []byte
}
