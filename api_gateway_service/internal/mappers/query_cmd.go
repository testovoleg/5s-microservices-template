package mappers

import (
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/app/commands"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
)

func NewAddApiCommand(accessToken, companyUuid string, in *dto.AddApiReqDto) *commands.AddApiCommand {
	return &commands.AddApiCommand{
		Params: &dto.ApiParamsDto{
			AccessToken: accessToken, CompanyUuid: companyUuid,
		},
		Dto: in,
	}
}

func NewGetApiCommand(accessToken string, in *dto.GetApiReqDto) *commands.GetApiCommand {
	return &commands.GetApiCommand{
		Params: &dto.ApiParamsDto{
			AccessToken: accessToken,
			CompanyUuid: in.CompanyUuid,
		},
	}
}

func NewGetFullApiCommand(accessToken string, in *dto.GetFullApiReqDto) *commands.GetFullApiCommand {
	return &commands.GetFullApiCommand{
		Params: &dto.ApiParamsDto{
			AccessToken: accessToken,
			CompanyUuid: in.CompanyUuid,
			ApiUuid:     in.ApiUuid,
		},
	}
}

func NewUpdateApiCommand(accessToken, companyUuid, apiUuid string, in *dto.UpdateApiReqDto) *commands.UpdateApiCommand {
	return &commands.UpdateApiCommand{
		Params: &dto.ApiParamsDto{
			AccessToken: accessToken,
			CompanyUuid: companyUuid,
			ApiUuid:     apiUuid,
		},
		Dto: in,
	}
}

func NewDeleteApiCommand(accessToken string, in *dto.DeleteApiReqDto) *commands.DeleteApiCommand {
	return &commands.DeleteApiCommand{
		Params: &dto.ApiParamsDto{
			AccessToken: accessToken,
			CompanyUuid: in.CompanyUuid,
			ApiUuid:     in.ApiUuid,
		},
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
