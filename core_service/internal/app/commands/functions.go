package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
	"go.opentelemetry.io/otel/trace"
)

// attrErr writes safe request params and payload into the span on error.
func attrErr(span trace.Span, params *models.ApiParams, data any, err error) error {
	if err == nil {
		return err
	}
	if params != nil {
		safeParams := *params
		safeParams.AccessToken = ""
		utils.AttrJSON(span, "params", safeParams)
	}
	utils.AttrJSON(span, "data", data)
	return httpErrors.ErrorWrap(err, span)
}

func getAdminToken(
	ctx context.Context,
	cloakRepo repository.IDMRepository, redisRepo repository.CacheRepository,
) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "commands.getAdminToken")
	defer span.End()

	adminToken, err := redisRepo.GetAdminToken(ctx)
	if err != nil {
		adminToken, err = cloakRepo.GetAdminToken(ctx)
		if err != nil {
			return "", httpErrors.ErrorWrap(err, span, "cloakRepo.GetAdminToken")
		}

		if adminToken == "" {
			return "", httpErrors.ErrorWrap(errors.New("can't get admin token from auth api"), span)
		}

		redisRepo.PutAdminToken(ctx, adminToken)

		utils.Attr(span, "admin_token_from", "auth_api")
		return adminToken, nil
	}

	utils.Attr(span, "admin_token_from", "redis")
	return adminToken, nil
}

// getCompanyApiData loads company database connection settings (ApiData) from Redis or Admin API.
func getCompanyApiData(
	ctx context.Context,
	cloakRepo repository.IDMRepository, adminRepo repository.AdminRepository, redisRepo repository.CacheRepository,
	company *models.Company,
) error {
	ctx, span := tracing.StartSpan(ctx, "commands.getCompanyApiData")
	defer span.End()

	if company.Uuid == "" {
		return httpErrors.ErrorWrap(errors.New("invalid input"), span)
	}

	adminToken, err := getAdminToken(ctx, cloakRepo, redisRepo)
	if err != nil {
		return err
	}

	apiData, err := redisRepo.GetApiData(ctx, company)
	if err == nil && apiData != nil {
		company.ApiData = apiData
		utils.Attr(span, "api_data_from", "redis")
	} else {
		err = adminRepo.GetApiData(ctx, adminToken, company)
		if err != nil {
			return httpErrors.ErrorWrap(err, span, "adminRepo.GetApiData")
		}
		_ = redisRepo.PutApiData(ctx, company, company.ApiData)
		utils.Attr(span, "api_data_from", "admin_api")
	}

	if company.ApiData == nil || company.ApiData.ApiUrl == "" || company.ApiData.Username == "" {
		return httpErrors.ErrorWrap(fmt.Errorf("something went wrong, check %s database connection settings", company.Uuid), span)
	}

	return nil
}

// getApiData resolves an API integration entity by uuid (Redis-first, then Admin property list).
func getApiData(
	ctx context.Context,
	adminRepo repository.AdminRepository, redisRepo repository.CacheRepository,
	params *models.ApiParams,
) (*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "commands.getApiData")
	defer span.End()

	if params == nil || params.AccessToken == "" || params.CompanyUuid == "" || params.ApiUuid == "" {
		return nil, httpErrors.ErrorWrap(errors.New("invalid input"), span)
	}

	utils.Attr(span, "api_uuid", params.ApiUuid)

	api, err := redisRepo.GetApi(ctx, params.CompanyUuid, params.ApiUuid)
	if err == nil && api != nil {
		return api, nil
	}

	listApi, err := adminRepo.GetProperty(ctx, params.AccessToken, params.CompanyUuid)
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "adminRepo.GetProperty")
	}

	redisRepo.PutApiList(ctx, params.CompanyUuid, listApi)

	for _, a := range listApi {
		if a.Uuid == params.ApiUuid {
			return a, nil
		}
	}

	return nil, httpErrors.ErrorWrap(errors.New("api with this uuid not found"), span)
}

func getUserData(
	ctx context.Context,
	log logger.Logger, cloakRepo repository.IDMRepository, adminRepo repository.AdminRepository, redisRepo repository.CacheRepository,
	params *models.ApiParams,
) (*models.User, *models.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "commands.getUserData")
	defer span.End()

	if params == nil {
		return nil, nil, httpErrors.ErrorWrap(errors.New("invalid input"), span)
	}

	u, err := cloakRepo.UserData(ctx, params.AccessToken)
	if err != nil {
		return nil, nil, httpErrors.ErrorWrap(err, span, "cloakRepo.UserData")
	}
	if u == nil {
		return nil, nil, httpErrors.ErrorWrap(errors.New("can't get user by token"), span)
	}
	if u.Company == nil {
		return nil, nil, httpErrors.ErrorWrap(errors.New("user has no associated company"), span)
	}

	utils.Attr(span, "token_idm_user_uuid", u.Id)
	utils.Attr(span, "token_idm_user_name", u.Name.FullName)
	utils.Attr(span, "token_idm_user_email", strings.Join(u.Contacts.Email, ";"))

	var (
		user    *models.User
		company *models.Company
	)

	isAdmin := cloakRepo.IsAdministrator(u)
	IsWebservice := cloakRepo.IsWebservice(u)
	IsSuperuser := cloakRepo.IsSuperuser(u)

	// Uncomment in a microservice that must allow external users without Admin/Webservice roles.
	// if cloakRepo.IsExternalUser(u) {
	// 	company = &models.Company{Uuid: u.Company.Uuid}
	// 	if params.CompanyUuid != "" {
	// 		company.Uuid = params.CompanyUuid
	// 	}
	// 	user = u
	// 	utils.Attr(span, "company_uuid", company.Uuid)
	// 	utils.Attr(span, "external_user", true)
	// 	if tempCo, _ := redisRepo.GetCompany(ctx, company.Uuid); tempCo != nil && tempCo.Uuid == company.Uuid {
	// 		utils.Attr(span, "company_name", tempCo.Name)
	// 		return user, tempCo, nil
	// 	}
	// 	return user, company, nil
	// }

	// delete this block if needed
	if !isAdmin && !IsSuperuser && !IsWebservice {
		return nil, nil, httpErrors.ErrorWrap(errors.New("access denied: this api methods requires user admin or webservice role"), span)
	}

	company = &models.Company{}
	if isAdmin || IsSuperuser {
		if params.CompanyUuid == "" {
			return nil, nil, httpErrors.ErrorWrap(errors.New("company_uuid is required for administrators"), span)
		}
		company.Uuid = params.CompanyUuid
	} else {
		company.Uuid = u.Company.Uuid
	}

	if u.Company.Uuid != company.Uuid && !isAdmin && !IsSuperuser {
		return nil, nil, httpErrors.ErrorWrap(errors.New("interaction with a foreign company is possible only for administrators"), span)
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
		return nil, nil, httpErrors.ErrorWrap(errors.New("access denied: cannot access another user's data"), span)
	}

	if targetUserID == u.Id {
		user = u
	} else {
		user, err = adminRepo.GetUserData(ctx, params.AccessToken, targetUserID)
		if err != nil {
			return nil, nil, httpErrors.ErrorWrap(errors.Wrapf(err, "adminRepo.GetUserData for user %s", targetUserID), span)
		}
		if user.Company == nil {
			return nil, nil, httpErrors.ErrorWrap(errors.New("target user has no associated company"), span)
		}
		if user.Company.Uuid != company.Uuid {
			return nil, nil, httpErrors.ErrorWrap(errors.New("user's company does not match the target company"), span)
		}
	}

	if tempCo, _ := redisRepo.GetCompany(ctx, company.Uuid); tempCo != nil && tempCo.Uuid == company.Uuid {
		utils.Attr(span, "company_name", tempCo.Name)
		return user, tempCo, nil
	}

	if adminToken, _ := getAdminToken(ctx, cloakRepo, redisRepo); adminToken != "" {
		if tempCo, _ := adminRepo.GetCompany(ctx, adminToken, company.Uuid); tempCo != nil && tempCo.Uuid == company.Uuid {
			redisRepo.PutCompany(ctx, tempCo)
			utils.Attr(span, "company_name", tempCo.Name)
			return user, tempCo, nil
		}
	}

	return user, company, nil
}
