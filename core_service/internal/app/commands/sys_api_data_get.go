package commands

import (
	"context"

	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

func getApiData(
	ctx context.Context,
	log logger.Logger, adminRepo repository.AdminRepository, redisRepo repository.CacheRepository,
	params *models.ApiParams,
) (*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "getApiCmdHandler.Handle")
	defer span.End()

	if params == nil || params.AccessToken == "" || params.CompanyUuid == "" || params.ApiUuid == "" {
		return nil, errors.New("invalid input")
	}

	utils.Attr(span, "api_uuid", params.ApiUuid)

	api, err := redisRepo.GetApi(ctx, params.CompanyUuid, params.ApiUuid)
	if err == nil && api != nil {
		return api, nil
	}

	listApi, err := adminRepo.GetProperty(ctx, params.AccessToken, params.CompanyUuid)
	if err != nil {
		return nil, errors.Wrap(err, "adminRepo.GetProperty")
	}

	redisRepo.PutApiList(ctx, params.CompanyUuid, listApi)

	for _, a := range listApi {
		if a.Uuid == params.ApiUuid {
			return a, nil
		}
	}

	return nil, errors.New("api with this uuid not found")
}
