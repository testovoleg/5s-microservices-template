package models

import (
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
)

func ListApiToGrpcMessage(in []*Api) []*coreService.Api {
	var list []*coreService.Api
	for _, v := range in {
		list = append(list, ApiToGrpcMessage(v))
	}
	return list
}

func ApiToGrpcMessage(in *Api) *coreService.Api {
	if in == nil {
		return &coreService.Api{}
	}

	res := &coreService.Api{
		Uuid:        in.Uuid,
		Title:       in.Title,
		Description: in.Description,
	}
	if in.Token != "" {
		res.Token = true
	}

	return res
}

func ApiFullToGrpcMessage(in *ApiFull) *coreService.ApiFull {
	if in == nil {
		return &coreService.ApiFull{}
	}

	res := &coreService.ApiFull{
		Uuid:        in.Uuid,
		Title:       in.Title,
		Description: in.Description,
		Token:       in.Token,
	}

	return res
}

func PresignUrlToGrpcMessage(in *PresignUrl) *coreService.PresignUrlRes {
	if in == nil {
		return nil
	}

	return &coreService.PresignUrlRes{
		FileId:      in.ObjectId,
		Url:         in.Url,
		Method:      in.Method,
		XAmxTagging: in.XAmzTagging,
	}
}
