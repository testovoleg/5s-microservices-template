package commands

import (
	"context"

	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

func getAdminToken(
	ctx context.Context,
	cloakRepo repository.IDMRepository, redisRepo repository.CacheRepository,
) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "getApiCmdHandler.Handle")
	defer span.End()

	adminToken, err := redisRepo.GetAdminToken(ctx)
	if err != nil {
		adminToken, err = cloakRepo.GetAdminToken(ctx)
		if err != nil {
			return "", errors.Wrap(err, "cloakRepo.GetAdminToken")
		}

		if adminToken == "" {
			return "", errors.New("can't get admin token from cloakRepo")
		}

		redisRepo.PutAdminToken(ctx, adminToken)

		utils.Attr(span, "admin_token_from", "auth_api")
		return adminToken, nil
	}

	utils.Attr(span, "admin_token_from", "redis")
	return adminToken, nil
}
