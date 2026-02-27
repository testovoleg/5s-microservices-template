package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
	"go.opentelemetry.io/otel/trace"
)

type adminRepository struct {
	log logger.Logger
	cfg *config.Config
}

func NewAdminRepository(log logger.Logger, cfg *config.Config) *adminRepository {
	return &adminRepository{log: log, cfg: cfg}
}

type PropertyData struct {
	List []*models.Api `json:"api_list"`
}

func (p *adminRepository) AddProperty(ctx context.Context, accessToken, companyUuid string, apiList []*models.Api) error {
	var span trace.Span
	ctx, span = tracing.StartSpan(ctx, "adminRepository.AddProperty")
	defer span.End()

	client := http.Client{
		Timeout: constants.HttpClientTimeout * time.Second,
	}
	requestURL := p.cfg.API.AdminApiUrl + "/property/" + p.cfg.API.PropertyKeyAPIData + "/company/" + companyUuid

	jsonData, err := json.Marshal(PropertyData{List: apiList})
	if err != nil {
		return err
	}

	body := map[string]interface{}{}
	body["value"] = string(jsonData)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	tracingHeaders, _ := tracing.InjectTextMapCarrier(ctx)
	for i, k := range tracingHeaders {
		req.Header.Set(i, k)
	}

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != http.StatusOK {
		return utils.NewError(res, span, "Error was during on Admin API")
	}

	return nil
}

func (p *adminRepository) GetProperty(ctx context.Context, accessToken, companyUuid string) ([]*models.Api, error) {
	var span trace.Span
	ctx, span = tracing.StartSpan(ctx, "adminRepository.GetProperty")
	defer span.End()

	client := http.Client{
		Timeout: constants.HttpClientTimeout * time.Second,
	}

	requestURL := p.cfg.API.AdminApiUrl + "/property/" + p.cfg.API.PropertyKeyAPIData + "/company/" + companyUuid

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	tracingHeaders, _ := tracing.InjectTextMapCarrier(ctx)
	for i, k := range tracingHeaders {
		req.Header.Set(i, k)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	if res.StatusCode != http.StatusOK {
		if strings.Contains(string(b), "no rows in result set") {
			return []*models.Api{}, nil
		}
		return nil, utils.NewError(res, span, "Error was during on Admin API")
	}

	var data PropertyData

	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	return data.List, nil
}

func (p *adminRepository) GetUserData(ctx context.Context, accessToken, userUuid string) (*models.User, error) {
	var span trace.Span
	ctx, span = tracing.StartSpan(ctx, "adminRepository.GetUserData")
	defer span.End()

	client := http.Client{
		Timeout: constants.HttpClientTimeout * time.Second,
	}

	requestURL := p.cfg.API.AdminApiUrl + "/user/" + userUuid

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	tracingHeaders, _ := tracing.InjectTextMapCarrier(ctx)
	for i, k := range tracingHeaders {
		req.Header.Set(i, k)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != http.StatusOK {
		return nil, utils.NewError(res, span, "Error was during on Admin API")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	type ResponseBody struct {
		Id      string          `json:"uuid"`
		Company *models.Company `json:"company"`
		Active  bool            `json:"active"`
		Roles   []string        `json:"roles"`
		Name    string          `json:"username"`
	}

	var bodyRes ResponseBody

	err = json.Unmarshal(b, &bodyRes)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	var roles []models.UserRole
	for _, r := range bodyRes.Roles {
		roles = append(roles, models.UserRole{Name: r})
	}

	user := &models.User{
		Id:      bodyRes.Id,
		Company: bodyRes.Company,
		Active:  bodyRes.Active,
		Roles:   roles,
		Name:    models.UserName{FullName: bodyRes.Name},
	}
	return user, nil
}

func (p *adminRepository) GetApi(ctx context.Context, accessToken, companyUuid, apiUuid string) (*models.Api, error) {
	listApi, err := p.GetProperty(ctx, accessToken, companyUuid)
	if err != nil {
		return nil, errors.Wrap(err, "adminRepo.GetApi")
	}

	for _, api := range listApi {
		if api.Uuid == apiUuid {
			return api, nil
		}
	}

	return nil, errors.New("api with this uuid not found")
}

func (r *adminRepository) GetCompany(ctx context.Context, accessToken, companyUuid string) (*models.Company, error) {
	var span trace.Span
	ctx, span = tracing.StartSpan(ctx, "adminRepository.GetCompany")
	defer span.End()
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestURL := r.cfg.API.AdminApiUrl + "/company/" + companyUuid

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	tracingHeaders, _ := tracing.InjectTextMapCarrier(ctx)
	for i, k := range tracingHeaders {
		req.Header.Set(i, k)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	type ResponeBody struct {
		Uuid string `json:"uuid"`
		Name string `json:"name"`
	}
	var resBody ResponeBody

	err = json.Unmarshal(b, &resBody)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	if res.StatusCode != http.StatusOK {
		return nil, utils.NewError(res, span, "Error was during on Admin API")
	}

	return &models.Company{
		Uuid: resBody.Uuid,
		Name: resBody.Name,
	}, nil
}
