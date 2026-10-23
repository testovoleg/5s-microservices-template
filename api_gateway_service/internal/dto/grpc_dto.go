package dto

import (
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

func ListApiDtoFromGrpc(in []*coreService.Api) []*ApiDto {
	var list []*ApiDto
	for _, v := range in {
		list = append(list, ApiDtoFromGrpc(v))
	}
	return utils.NonNilSlice(list)
}

func ApiDtoFromGrpc(in *coreService.Api) *ApiDto {
	if in == nil {
		return &ApiDto{}
	}
	return &ApiDto{
		Uuid:        in.GetUuid(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Token:       in.GetToken(),
	}
}

func ApiFullDtoFromGrpc(in *coreService.ApiFull) *ApiFullDto {
	if in == nil {
		return &ApiFullDto{Subscriptions: []string{}}
	}

	return &ApiFullDto{
		Uuid:          in.GetUuid(),
		Title:         in.GetTitle(),
		Description:   in.GetDescription(),
		Token:         in.GetToken(),
		Subscriptions: utils.NonNilSlice(in.GetSubscriptions()),
	}
}

func EventDtoFromGrpc(in *coreService.Event) *EventDto {
	if in == nil {
		return nil
	}
	return &EventDto{
		Uuid:        in.GetUuid(),
		Description: in.GetDescription(),
		Timestamp:   in.GetTimestamp(),
		Status:      EventStatusDto(in.GetStatus()),
	}
}
