package repository

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

const (
	redisMicroservicePrefixKey = constants.ShortMicroserviceName
	redisAdminTokenLifetime    = 3 * time.Hour
	redisApiLife               = time.Hour * 24 * 30
	redisEventLife             = time.Hour * 24
	// redisCompanyFullListLife — короткий TTL кэша полного списка компаний: справочник
	// меняется редко, но новая компания должна появиться в админских выборках быстро.
	redisCompanyFullListLife = 5 * time.Minute
)

type redisRepository struct {
	log         logger.Logger
	cfg         *config.Config
	redisClient redis.UniversalClient
}

func NewRedisRepository(log logger.Logger, cfg *config.Config, redisClient redis.UniversalClient) *redisRepository {
	return &redisRepository{log: log, cfg: cfg, redisClient: redisClient}
}

func (r *redisRepository) GetApiData(ctx context.Context, company *models.Company) (*models.ApiData, error) {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.GetApiData")
	defer span.End()

	key := r.getRedisApiDataListKey()
	field := redisGlobEscape(company.Uuid)

	fields, err := r.HScan(ctx, key, field)
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span)
	}

	for _, f := range fields {
		var item models.ApiData
		if err := json.Unmarshal([]byte(f.Value), &item); err != nil {
			continue
		}
		return &item, nil
	}

	return nil, nil
}

func (r *redisRepository) PutApiData(ctx context.Context, company *models.Company, apiData *models.ApiData) error {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.PutApiData")
	defer span.End()

	if company == nil || company.Uuid == "" {
		return httpErrors.ErrorWrap(errors.New("invalid company. Nil or empty company uuid"), span)
	}

	if apiData == nil || apiData.ApiUrl == "" {
		return httpErrors.ErrorWrap(errors.New("invalid api data. Nil or empty api_url"), span)
	}

	key := r.getRedisApiDataListKey()
	apiBytes, err := json.Marshal(apiData)
	if err != nil {
		return httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	field := company.Uuid

	if err := r.redisClient.HSet(ctx, key, field, apiBytes).Err(); err != nil {
		return httpErrors.ErrorWrap(err, span, "redisClient.HSet")
	}
	r.redisClient.HExpire(ctx, key, 24*time.Hour, field)

	return nil
}

func (r *redisRepository) PutAdminToken(ctx context.Context, token string) error {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.PutAdminToken")
	defer span.End()

	if token == "" {
		return httpErrors.ErrorWrap(errors.New("invalid input"), span)
	}

	if err := r.redisClient.Set(ctx, r.getRedisAdminTokenKey(), token, redisAdminTokenLifetime).Err(); err != nil {
		return httpErrors.ErrorWrap(err, span, "redisClient.Set")
	}

	return nil
}

func (r *redisRepository) GetAdminToken(ctx context.Context) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.GetAdminToken")
	defer span.End()

	if token := r.redisClient.Get(ctx, r.getRedisAdminTokenKey()).Val(); token != "" {
		return token, nil
	}

	return "", httpErrors.ErrorWrap(errors.New("redisClient.Get"), span)
}

func (r *redisRepository) PutApiList(ctx context.Context, companyUuid string, apiList []*models.Api) error {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.PutApiList")
	defer span.End()

	key := r.getRedisApiListKey()

	for _, api := range apiList {
		apiBytes, err := json.Marshal(api)
		if err != nil {
			r.log.WarnMsg("json.Marshal", err)
			return httpErrors.ErrorWrap(err, span, "json.Marshal")
		}

		field := r.getRedisApiListField(companyUuid, api.Uuid)

		if err := r.redisClient.HSet(ctx, key, field, apiBytes).Err(); err != nil {
			r.log.WarnMsg("redisClient.HSet", err)
			return httpErrors.ErrorWrap(err, span, "redisClient.HSet")
		}
		r.redisClient.HExpire(ctx, key, redisApiLife, field)
	}
	r.redisClient.Expire(ctx, key, redisApiLife)

	return nil
}

func (r *redisRepository) GetApiList(ctx context.Context, companyUuid string) ([]*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.GetApiList")
	defer span.End()

	key := r.getRedisApiListKey()
	field := r.getRedisApiListMatch(companyUuid, constants.RedisAllItems)

	fields, err := r.HScan(ctx, key, field)
	if err != nil {
		return nil, err
	}

	var res []*models.Api
	for _, f := range fields {
		var item models.Api
		if err := json.Unmarshal([]byte(f.Value), &item); err != nil {
			continue
		}
		res = append(res, &item)
	}

	return res, nil
}

func (r *redisRepository) GetApi(ctx context.Context, companyUuid, apiUuid string) (*models.Api, error) {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.GetApi")
	defer span.End()

	if apiUuid == "" {
		return nil, errors.New("empty input")
	}

	key := r.getRedisApiListKey()
	field := r.getRedisApiListMatch(companyUuid, apiUuid)

	fields, err := r.HScan(ctx, key, field)
	if err != nil {
		return nil, err
	}

	for _, f := range fields {
		var item models.Api
		if err := json.Unmarshal([]byte(f.Value), &item); err != nil {
			continue
		}
		return &item, nil
	}

	return nil, nil
}

func (r *redisRepository) DeleteApiList(ctx context.Context, company *models.Company) error {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.DeleteApiList")
	defer span.End()

	key := r.getRedisApiListKey()
	field := r.getRedisApiListMatch(company.Uuid, constants.RedisAllItems)

	fields, err := r.HScan(ctx, key, field)
	if err != nil {
		return err
	}

	for _, f := range fields {
		r.redisClient.HDel(ctx, key, f.Field)
	}

	return nil
}

func (r *redisRepository) GetCompany(ctx context.Context, companyUuid string) (*models.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.GetCompany")
	defer span.End()

	key := r.getRedisCompanyListKey()
	field := redisGlobEscape(companyUuid)

	fields, err := r.HScan(ctx, key, field)
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span)
	}

	for _, f := range fields {
		var item models.Company
		if err := json.Unmarshal([]byte(f.Value), &item); err != nil {
			continue
		}
		return &item, nil
	}

	return nil, nil
}

func (r *redisRepository) PutCompany(ctx context.Context, company *models.Company) error {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.PutCompany")
	defer span.End()

	if company == nil || company.Uuid == "" || company.Name == "" {
		return httpErrors.ErrorWrap(errors.New("invalid company. Nil or empty company uuid"), span)
	}

	key := r.getRedisCompanyListKey()
	bytes, err := json.Marshal(company)
	if err != nil {
		return httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	field := company.Uuid

	if err := r.redisClient.HSet(ctx, key, field, bytes).Err(); err != nil {
		return httpErrors.ErrorWrap(err, span, "redisClient.HSet")
	}
	r.redisClient.HExpire(ctx, key, 24*time.Hour, field)

	return nil
}

type HashField struct {
	Field string
	Value string
}

func (r *redisRepository) HScan(ctx context.Context, key, matchPattern string) ([]*HashField, error) {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.HScan")
	defer span.End()

	var results []*HashField
	var cursor uint64

	for {
		fields, nextCursor, err := r.redisClient.HScan(ctx, key, cursor, matchPattern, 0).Result()
		if err != nil {
			return nil, httpErrors.ErrorWrap(err, span, "redisClient.HScan")
		}

		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				results = append(results, &HashField{
					Field: fields[i],
					Value: fields[i+1],
				})
			}
		}

		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	return results, nil
}

func (r *redisRepository) getRedisAdminTokenKey() string {
	return r.getRedisMicroservicePrefixKey() + ":admin"
}

// GetCompanyList returns the cached full company directory (see PutCompanyList), or
// (nil, nil) on a cache miss/expiry — callers fall back to adminRepository.GetCompanyList.
func (r *redisRepository) GetCompanyList(ctx context.Context) ([]*models.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.GetCompanyList")
	defer span.End()

	val := r.redisClient.Get(ctx, r.getRedisCompanyFullListKey()).Val()
	if val == "" {
		return nil, nil
	}

	var companies []*models.Company
	if err := json.Unmarshal([]byte(val), &companies); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}
	return companies, nil
}

// PutCompanyList caches the full company directory fetched from the Admin API for
// redisCompanyFullListLife.
func (r *redisRepository) PutCompanyList(ctx context.Context, companies []*models.Company) error {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.PutCompanyList")
	defer span.End()

	bytes, err := json.Marshal(companies)
	if err != nil {
		return httpErrors.ErrorWrap(err, span, "json.Marshal")
	}
	if err := r.redisClient.Set(ctx, r.getRedisCompanyFullListKey(), bytes, redisCompanyFullListLife).Err(); err != nil {
		return httpErrors.ErrorWrap(err, span, "redisClient.Set")
	}
	return nil
}

func (r *redisRepository) getRedisCompanyFullListKey() string {
	return r.getRedisMicroservicePrefixKey() + ":company_full_list"
}

func (r *redisRepository) getRedisMicroservicePrefixKey() string {
	if r.cfg.ServiceSettings.RedisMicroservicePrefixKey != "" {
		return r.cfg.ServiceSettings.RedisMicroservicePrefixKey
	}

	return redisMicroservicePrefixKey
}

func (r *redisRepository) getRedisApiListKey() string {
	return r.getRedisMicroservicePrefixKey() + ":api_list"
}

func (r *redisRepository) getRedisApiListField(companyUuid, apiUuid string) string {
	return "company@" + companyUuid + ":api@" + apiUuid
}

func (r *redisRepository) getRedisApiListMatch(companyUuid, apiUuid string) string {
	company := redisGlobEscape(companyUuid)
	if apiUuid == constants.RedisAllItems {
		return "company@" + company + ":api@*"
	}
	return "company@" + company + ":api@" + redisGlobEscape(apiUuid)
}

// redisGlobEscape treats user-supplied ids as literals in Redis MATCH (so * ? [ ] cannot dump other keys).
func redisGlobEscape(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '*', '?', '[', ']', '\\':
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func (r *redisRepository) getRedisApiDataListKey() string {
	return r.getRedisMicroservicePrefixKey() + ":api_data_list"
}

func (r *redisRepository) getRedisCompanyListKey() string {
	return r.getRedisMicroservicePrefixKey() + ":company_list"
}
