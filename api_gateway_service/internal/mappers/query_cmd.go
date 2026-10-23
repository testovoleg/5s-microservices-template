package mappers

import (
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/app/commands"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
)

func NewAddApiCommand(params *dto.ApiParamsDto, in *dto.AddApiReqDto) *commands.AddApiCommand {
	return &commands.AddApiCommand{
		Params: params,
		Dto:    in,
	}
}

func NewGetApiCommand(params *dto.ApiParamsDto) *commands.GetApiCommand {
	return &commands.GetApiCommand{
		Params: params,
	}
}

func NewGetFullApiCommand(params *dto.ApiParamsDto) *commands.GetFullApiCommand {
	return &commands.GetFullApiCommand{
		Params: params,
	}
}

func NewUpdateApiCommand(params *dto.ApiParamsDto, in *dto.UpdateApiReqDto) *commands.UpdateApiCommand {
	return &commands.UpdateApiCommand{
		Params: params,
		Dto:    in,
	}
}

func NewDeleteApiCommand(params *dto.ApiParamsDto) *commands.DeleteApiCommand {
	return &commands.DeleteApiCommand{
		Params: params,
	}
}

func NewWebhookCommand(companyUuid, apiUuid string, in []byte) *commands.WebhookCommand {
	return &commands.WebhookCommand{
		Params: &dto.ApiParamsDto{
			CompanyUuid: companyUuid,
			ApiUuid:     apiUuid,
		},
		Payload: in,
	}
}

func NewEventGetCommand(params *dto.ApiParamsDto, in *dto.EventGetReqDto) *commands.EventGetCommand {
	return &commands.EventGetCommand{
		Params: params,
		Dto:    in,
	}
}
