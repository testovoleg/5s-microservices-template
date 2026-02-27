package commands

import kafkaMessages "github.com/testovoleg/5s-microservice-template/proto/kafka"

type WebhookCommand struct {
	Payload *kafkaMessages.Payload
}
