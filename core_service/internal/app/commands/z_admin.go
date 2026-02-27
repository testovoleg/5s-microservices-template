package commands

import (
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
)

type AddApiCommand struct {
	Params *models.ApiParams
	Data   *models.Api
}

type GetApiCommand struct {
	Params *models.ApiParams
}

type GetFullApiCommand struct {
	Params *models.ApiParams
}

type UpdateApiCommand struct {
	Params *models.ApiParams
	Data   *models.UpdateApi
}

type DeleteApiCommand struct {
	Params *models.ApiParams
}
type SysGetApiCommand struct {
	Params *models.ApiParams
}

type SysGetUserDataCommand struct {
	IdmUserUuid string `json:"idm_user_uuid"`
	Params      *models.ApiParams
}

type SysGetFileDataCommand struct {
	FileId string `json:"file_id"`
	Params *models.ApiParams
}
