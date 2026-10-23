package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	httpUtils "github.com/testovoleg/5s-microservice-template/pkg/http_utils"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

type adminRepository struct {
	log logger.Logger
	cfg *config.Config
}

func NewAdminRepository(log logger.Logger, cfg *config.Config) *adminRepository {
	return &adminRepository{log: log, cfg: cfg}
}

func (p *adminRepository) AddProperty(ctx context.Context, accessToken, companyUuid string, apiList []*models.Api) error {
	ctx, span := tracing.StartSpan(ctx, "adminRepository.AddProperty")
	defer span.End()

	requestURL := p.cfg.API.AdminApiUrl + "/property/" + p.cfg.API.PropertyKeyAPIData + "/company/" + companyUuid

	jsonData, err := json.Marshal(models.PropertyData{List: apiList})
	if err != nil {
		return err
	}

	jsonBody, err := json.Marshal(models.PropertyValueReq{Value: string(jsonData)})
	if err != nil {
		return err
	}

	_, err = httpUtils.SendHTTPRequest(ctx, http.MethodPost, requestURL, jsonBody, httpUtils.SendOpts{
		BearerToken: accessToken,
		ErrorPrefix: "Error was during on Admin API",
	})
	return err
}

func (p *adminRepository) GetProperty(ctx context.Context, accessToken, companyUuid string) ([]*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "adminRepository.GetProperty")
	defer span.End()

	requestURL := p.cfg.API.AdminApiUrl + "/property/" + p.cfg.API.PropertyKeyAPIData + "/company/" + companyUuid

	status, b, err := httpUtils.SendHTTPRequestRaw(ctx, http.MethodGet, requestURL, nil, httpUtils.SendOpts{
		BearerToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		if strings.Contains(string(b), "no rows in result set") {
			return []*models.Api{}, nil
		}
		return nil, httpErrors.HttpBody(fmt.Sprintf("%d %s", status, http.StatusText(status)), b, span, "Error was during on Admin API")
	}

	var data models.PropertyData
	if err = json.Unmarshal(b, &data); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	return data.List, nil
}

func (p *adminRepository) GetUserData(ctx context.Context, accessToken, userUuid string) (*models.User, error) {
	ctx, span := tracing.StartSpan(ctx, "adminRepository.GetUserData")
	defer span.End()

	requestURL := p.cfg.API.AdminApiUrl + "/user/" + userUuid

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodGet, requestURL, nil, httpUtils.SendOpts{
		BearerToken: accessToken,
		ContentType: "application/json",
		ErrorPrefix: "Error was during on Admin API",
	})
	if err != nil {
		return nil, err
	}

	var bodyRes models.AdminUserRes
	if err = json.Unmarshal(b, &bodyRes); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	var roles []models.UserRole
	for _, r := range bodyRes.Roles {
		roles = append(roles, models.UserRole{Name: r})
	}

	return &models.User{
		Id:      bodyRes.Id,
		Company: bodyRes.Company,
		Active:  bodyRes.Active,
		Roles:   roles,
		Name:    models.UserName{FullName: bodyRes.Name},
	}, nil
}

func (p *adminRepository) GetApi(ctx context.Context, accessToken, companyUuid, apiUuid string) (*models.Api, error) {
	listApi, err := p.GetProperty(ctx, accessToken, companyUuid)
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, nil, "adminRepo.GetApi")
	}

	for _, api := range listApi {
		if api.Uuid == apiUuid {
			return api, nil
		}
	}

	return nil, errors.New("api with this uuid not found")
}

func (r *adminRepository) GetCompany(ctx context.Context, accessToken, companyUuid string) (*models.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "adminRepository.GetCompany")
	defer span.End()

	requestURL := r.cfg.API.AdminApiUrl + "/company/" + companyUuid

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodGet, requestURL, nil, httpUtils.SendOpts{
		BearerToken: accessToken,
		Timeout:     30 * time.Second,
		ErrorPrefix: "Error was during on Admin API",
	})
	if err != nil {
		return nil, err
	}

	var resBody models.AdminCompanyRes
	if err = json.Unmarshal(b, &resBody); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	return &models.Company{
		Uuid:   resBody.Uuid,
		Name:   resBody.Name,
		Active: utils.DerefBool(resBody.Active),
	}, nil
}

func (r *adminRepository) GetCompanyList(ctx context.Context, accessToken string) ([]*models.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "adminRepository.GetCompanyList")
	defer span.End()

	// Admin /company/list is 1-based (same convention as pkg/utils.Pagination: page 1 is
	// offset 0, page 0 would be the same first page again). Walking from 0 duplicates
	// every company on page 1.
	seen := make(map[string]struct{})
	var out []*models.Company
	for page := 1; page <= models.AdminCompanyListMaxPages; page++ {
		requestURL := fmt.Sprintf("%s/company/list?page=%d&size=%d", r.cfg.API.AdminApiUrl, page, models.AdminCompanyListPageSize)

		b, err := httpUtils.SendHTTPRequest(ctx, http.MethodGet, requestURL, nil, httpUtils.SendOpts{
			BearerToken: accessToken,
			Timeout:     30 * time.Second,
			ErrorPrefix: "Error was during on Admin API",
		})
		if err != nil {
			return nil, err
		}

		var res models.AdminCompanyList
		if err = json.Unmarshal(b, &res); err != nil {
			return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
		}
		for _, c := range res.Companies {
			if c.Uuid == "" {
				continue
			}
			if _, dup := seen[c.Uuid]; dup {
				continue
			}
			seen[c.Uuid] = struct{}{}
			out = append(out, &models.Company{
				Uuid:   c.Uuid,
				Name:   c.Name,
				Active: utils.DerefBool(c.Active),
			})
		}
		if !res.HasMore || len(res.Companies) == 0 {
			break
		}
		if res.TotalPages > 0 && page >= res.TotalPages {
			break
		}
	}
	return out, nil
}

func (r *adminRepository) GetApiData(ctx context.Context, accessToken string, company *models.Company) (err error) {
	ctx, span := tracing.StartSpan(ctx, "adminRepository.GetApiData")
	defer span.End()

	requestURL := r.cfg.API.AdminApiUrl + "/property/" + constants.UniversalAPI + "/company/" + company.Uuid

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodGet, requestURL, nil, httpUtils.SendOpts{
		BearerToken: accessToken,
	})
	if err != nil {
		return err
	}

	var body models.ApiData
	if err = json.Unmarshal(b, &body); err != nil {
		return httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	company.ApiData = &body
	if company.ApiData == nil {
		company.ApiData = &models.ApiData{}
	}

	return nil
}
