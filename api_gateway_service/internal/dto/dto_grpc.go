package dto

import (
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
)

func ApiParamsToGrpcMessage(in *ApiParamsDto) *coreService.ApiParams {
	if in == nil {
		return &coreService.ApiParams{}
	}
	return &coreService.ApiParams{
		AccessToken: in.AccessToken,
		CompanyUuid: in.CompanyUuid,
		ApiUuid:     in.ApiUuid,
	}
}

func AddApiToGrpcMessage(p *ApiParamsDto, in *AddApiReqDto) *coreService.AddApiReq {
	return &coreService.AddApiReq{
		Params:      ApiParamsToGrpcMessage(p),
		Title:       in.Title,
		Description: in.Description,
		Token:       in.Token,
	}
}

func GetApiToGrpcMessage(p *ApiParamsDto) *coreService.GetApiReq {
	return &coreService.GetApiReq{
		Params: ApiParamsToGrpcMessage(p),
	}
}

func GetFullApiToGrpcMessage(p *ApiParamsDto) *coreService.GetFullApiReq {
	return &coreService.GetFullApiReq{
		Params: ApiParamsToGrpcMessage(p),
	}
}

func UpdateApiToGrpcMessage(p *ApiParamsDto, in *UpdateApiReqDto) *coreService.UpdateApiReq {
	return &coreService.UpdateApiReq{
		Params:      ApiParamsToGrpcMessage(p),
		Title:       in.Title,
		Description: in.Description,
		Token:       in.Token,
	}
}

func DeleteApiToGrpcMessage(p *ApiParamsDto) *coreService.DeleteApiReq {
	return &coreService.DeleteApiReq{
		Params: ApiParamsToGrpcMessage(p),
	}
}
