package model

import (
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

func GeneralParamsInputToGrpc(in *GeneralParamsInput) *coreService.ApiParams {
	if in == nil {
		return &coreService.ApiParams{}
	}
	return &coreService.ApiParams{
		AccessToken: utils.DerefString(in.AccessToken),
		CompanyUuid: in.CompanyUUID,
		ApiUuid:     utils.DerefString(in.APIUUID),
		IdmUserUuid: utils.DerefString(in.IdmUserUUID),
	}
}

func EventGetInputToGrpc(in *TipicalEventGetInput) *coreService.EventGetReq {
	if in == nil {
		return &coreService.EventGetReq{}
	}
	return &coreService.EventGetReq{
		EventUuid: in.EventUUID,
	}
}

func EventFromGrpc(in *coreService.Event) *TipicalEvent {
	if in == nil {
		return nil
	}
	return &TipicalEvent{
		UUID:        utils.StrPtr(in.GetUuid()),
		Description: utils.StrPtr(in.GetDescription()),
		Timestamp:   utils.StrPtr(in.GetTimestamp()),
		Status:      utils.StrPtr(in.GetStatus()),
	}
}
