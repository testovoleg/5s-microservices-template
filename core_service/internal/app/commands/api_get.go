package commands

import (
	"context"

	"github.com/pkg/errors"

	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

func (c *apiMethodsHandler) GetApi(ctx context.Context, command *GetApiCommand) ([]*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "apiMethodsHandler.GetApi")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, c.redisRepo, command.Params)
	if err != nil {
		return nil, attrErr(span, command.Params, nil, err)
	}

	listApi, err := c.redisRepo.GetApiList(ctx, company.Uuid)
	if err != nil || len(listApi) == 0 {
		listApi, err = c.adminRepo.GetProperty(ctx, command.Params.AccessToken, company.Uuid)
		if err != nil {
			return nil, attrErr(span, command.Params, nil, errors.Wrap(err, "adminRepo.GetProperty"))
		}
		c.redisRepo.PutApiList(ctx, company.Uuid, listApi)
	}

	return listApi, nil
}
