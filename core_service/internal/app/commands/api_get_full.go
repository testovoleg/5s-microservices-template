package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

func (c *apiMethodsHandler) GetFullApi(ctx context.Context, command *GetFullApiCommand) (*models.ApiFull, error) {
	ctx, span := tracing.StartSpan(ctx, "apiMethodsHandler.GetFullApi")
	defer span.End()

	_, company, err := getUserData(ctx, c.log, c.cloakRepo, c.adminRepo, command.Params)
	if err != nil {
		return nil, err
	}

	api, err := getApiData(ctx, c.adminRepo, c.redisRepo, &models.ApiParams{AccessToken: command.Params.AccessToken, CompanyUuid: company.Uuid, ApiUuid: command.Params.ApiUuid})
	if err != nil {
		return nil, err
	}

	subscriptions := []string{}
	// subscriptions, err := c.targetRepo.GetSubscriptions(ctx, company, api)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "targetRepo.GetSubscriptions")
	// }

	return &models.ApiFull{
		Uuid:          api.Uuid,
		Title:         api.Title,
		Description:   api.Description,
		Token:         api.Token,
		Subscriptions: subscriptions,
	}, nil
}
