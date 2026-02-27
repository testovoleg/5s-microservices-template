package dto

import (
	kafkaMessages "github.com/testovoleg/5s-microservice-template/proto/kafka"
)

func ApiParamsDtoToKafkaMessage(in *ApiParamsDto) *kafkaMessages.ApiParams {
	if in == nil {
		return nil
	}

	return &kafkaMessages.ApiParams{
		AccessToken: in.AccessToken,
		CompanyUuid: in.CompanyUuid,
		ApiUuid:     in.ApiUuid,
	}
}
