package mappers

import (
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/commands"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
)

func NewAddApiCommandFromGrpcMessage(msg *coreService.AddApiReq) *commands.AddApiCommand {
	return &commands.AddApiCommand{
		Params: models.ApiParamsFromGrpcMessage(msg.GetParams()),
		Data: &models.Api{
			Title:       msg.GetTitle(),
			Description: msg.GetDescription(),
			Token:       msg.GetToken(),
		},
	}
}

func NewGetApiCommandFromGrpcMessage(msg *coreService.GetApiReq) *commands.GetApiCommand {
	return &commands.GetApiCommand{
		Params: models.ApiParamsFromGrpcMessage(msg.GetParams()),
	}
}

func NewGetFullApiCommandFromGrpcMessage(msg *coreService.GetFullApiReq) *commands.GetFullApiCommand {
	return &commands.GetFullApiCommand{
		Params: models.ApiParamsFromGrpcMessage(msg.GetParams()),
	}
}

func NewUpdateApiCommandFromGrpcMessage(msg *coreService.UpdateApiReq) *commands.UpdateApiCommand {
	return &commands.UpdateApiCommand{
		Params: models.ApiParamsFromGrpcMessage(msg.GetParams()),
		Data: &models.UpdateApi{
			Title:       msg.Title,
			Description: msg.Description,
			Token:       msg.Token,
		},
	}
}

func NewDeleteApiCommandFromGrpcMessage(msg *coreService.DeleteApiReq) *commands.DeleteApiCommand {
	return &commands.DeleteApiCommand{
		Params: models.ApiParamsFromGrpcMessage(msg.GetParams()),
	}
}
