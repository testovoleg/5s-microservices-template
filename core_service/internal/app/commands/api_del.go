package commands

import (
	"context"

	"github.com/pkg/errors"

	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

func (c *apiMethodsHandler) DeleteApi(ctx context.Context, command *DeleteApiCommand) error {
	ctx, span := tracing.StartSpan(ctx, "apiMethodsHandler.DeleteApi")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, command.Params)
	if err != nil {
		return err
	}

	if command.Params.ApiUuid == "" {
		return errors.New("apiUuid is not filled")
	}

	listApi, err := c.adminRepo.GetProperty(ctx, command.Params.AccessToken, company.Uuid)
	if err != nil {
		return errors.Wrap(err, "adminRepo.GetProperty")
	}

	res, deleted := c.filterApiList(listApi, command.Params.ApiUuid)
	if !deleted {
		return errors.New("apiUuid not found")
	}

	//delete before add if admin change with error
	c.redisRepo.DeleteApiList(ctx, company)

	err = c.adminRepo.AddProperty(ctx, command.Params.AccessToken, company.Uuid, res)
	if err != nil {
		return errors.Wrap(err, "adminRepo.AddProperty")
	}

	return nil
}

func (c *apiMethodsHandler) filterApiList(list []*models.Api, apiUuid string) (res []*models.Api, deleted bool) {
	for _, api := range list {
		if api.Uuid == apiUuid {
			deleted = true
			continue
		}

		res = append(res, api)
	}
	return
}
