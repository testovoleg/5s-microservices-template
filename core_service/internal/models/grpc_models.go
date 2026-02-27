package models

import coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"

func ApiParamsFromGrpcMessage(in *coreService.ApiParams) *ApiParams {
	return &ApiParams{
		AccessToken: in.GetAccessToken(),
		CompanyUuid: in.GetCompanyUuid(),
		ApiUuid:     in.GetApiUuid(),
		IdmUserUuid: in.GetIdmUserUuid(),
	}
}
