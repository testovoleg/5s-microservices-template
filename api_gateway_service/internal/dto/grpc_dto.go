package dto

import (
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
)

func ListApiDtoFromGrpc(in []*coreService.Api) []*ApiDto {
	if in == nil {
		return []*ApiDto{}
	}

	var list []*ApiDto
	for _, v := range in {
		list = append(list, ApiDtoFromGrpc(v))
	}
	return list
}

func ApiDtoFromGrpc(in *coreService.Api) *ApiDto {
	return &ApiDto{
		Uuid:        in.GetUuid(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Token:       in.GetToken(),
	}
}

func ApiFullDtoFromGrpc(in *coreService.ApiFull) *ApiFullDto {
	if in == nil {
		return &ApiFullDto{}
	}

	return &ApiFullDto{
		Uuid:        in.GetUuid(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Token:       in.GetToken(),
	}
}
