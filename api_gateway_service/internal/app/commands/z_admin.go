package commands

import "github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"

type AddApiCommand struct {
	Params *dto.ApiParamsDto
	Dto    *dto.AddApiReqDto
}

type GetApiCommand struct {
	Params *dto.ApiParamsDto
}

type GetFullApiCommand struct {
	Params *dto.ApiParamsDto
}

type UpdateApiCommand struct {
	Params *dto.ApiParamsDto
	Dto    *dto.UpdateApiReqDto
}

type DeleteApiCommand struct {
	Params *dto.ApiParamsDto
}
