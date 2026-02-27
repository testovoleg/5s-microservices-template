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

func getUserData(
	ctx context.Context,
	log logger.Logger, cloakRepo repository.IDMRepository, adminRepo repository.AdminRepository, redisRepo repository.CacheRepository,
	params *models.ApiParams,
) (*models.User, *models.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "getUserDataCmdHandler.Handle")
	defer span.End()

	if params == nil {
		return nil, nil, errors.New("invalid input")
	}

	u, err := cloakRepo.UserData(ctx, params.AccessToken)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cloakRepo.UserData")
	}
	if u == nil {
		return nil, nil, errors.New("can't get user by token")
	}

	utils.Attr(span, "token_idm_user_uuid", u.Id)

	var (
		user    *models.User
		company *models.Company
	)

	isAdmin := cloakRepo.IsAdministrator(u)
	IsWebservice := cloakRepo.IsWebservice(u)
	IsSuperuser := cloakRepo.IsSuperuser(u)

	//delete this block if needed
	if !isAdmin && !IsSuperuser && !IsWebservice {
		return nil, nil, errors.New("access denied: this api methods requires user admin or webservice role")
	}
	//

	company = &models.Company{}
	if isAdmin || IsSuperuser {
		if params.CompanyUuid == "" {
			return nil, nil, errors.New("company_uuid is required for administrators")
		}
		company.Uuid = params.CompanyUuid
	} else {
		if u.Company == nil {
			return nil, nil, errors.New("user has no associated company")
		}
		company.Uuid = u.Company.Uuid
	}

	utils.Attr(span, "company_uuid", company.Uuid)

	targetUserID := params.IdmUserUuid
	if targetUserID == "" {
		targetUserID = u.Id
	}

	utils.Attr(span, "target_idm_user_uuid", targetUserID)

	if !isAdmin && !IsSuperuser && !IsWebservice && targetUserID != u.Id {
		log.Warn(ctx, "access denied: non-admin tried to access another user",
			"requester_id", u.Id, "target_id", targetUserID)
		return nil, nil, errors.New("access denied: cannot access another user's data")
	}

	if targetUserID == u.Id {
		user = u
	} else {
		user, err = adminRepo.GetUserData(ctx, params.AccessToken, targetUserID)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "adminRepo.GetUserData for user %s", targetUserID)
		}
		if user.Company == nil {
			return nil, nil, errors.New("target user has no associated company")
		}
		if user.Company.Uuid != company.Uuid {
			return nil, nil, errors.New("user's company does not match the target company")
		}
	}

	return user, company, nil
}
