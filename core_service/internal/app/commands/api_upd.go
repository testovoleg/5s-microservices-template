package commands

import (
	"context"

	"github.com/pkg/errors"

	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

func (c *apiMethodsHandler) UpdateApi(ctx context.Context, command *UpdateApiCommand) (*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "apiMethodsHandler.UpdateApi")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, c.redisRepo, command.Params)
	if err != nil {
		return nil, attrErr(span, command.Params, command.Data, err)
	}

	if command.Params.ApiUuid == "" {
		return nil, attrErr(span, command.Params, command.Data, errors.New("apiUuid is not filled"))
	}

	//do not use redis just in case
	listApi, err := c.adminRepo.GetProperty(ctx, command.Params.AccessToken, company.Uuid)
	if err != nil {
		return nil, attrErr(span, command.Params, command.Data, errors.Wrap(err, "adminRepo.GetProperty"))
	}

	api, updated := c.updateApi(listApi, command.Data, command.Params.ApiUuid)
	if !updated {
		return nil, attrErr(span, command.Params, command.Data, errors.New("apiUuid not found"))
	}

	err = c.adminRepo.AddProperty(ctx, command.Params.AccessToken, company.Uuid, listApi)
	if err != nil {
		return nil, attrErr(span, command.Params, command.Data, errors.Wrap(err, "adminRepo.AddProperty"))
	}

	c.redisRepo.PutApiList(ctx, company.Uuid, listApi)

	return api, nil
}

func (c *apiMethodsHandler) updateApi(list []*models.Api, update *models.UpdateApi, apiUuid string) (api *models.Api, updated bool) {
	for _, a := range list {
		if a == nil || a.Uuid != apiUuid {
			continue
		}
		a.Title = utils.CoalescePatchString(a.Title, update.Title, "string")
		a.Description = utils.CoalescePatchString(a.Description, update.Description, "string")
		a.Token = utils.CoalescePatchString(a.Token, update.Token, "string")
		api = a
		updated = true
	}

	return
}
