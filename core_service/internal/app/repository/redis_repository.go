package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

const (
	redisMicroservicePrefixKey = constants.ShortMicroserviceName
	redisAdminTokenLifetime    = 3 * time.Hour
	redisApiLife               = time.Hour * 24 * 30
)

type redisRepository struct {
	log         logger.Logger
	cfg         *config.Config
	redisClient redis.UniversalClient
}

func NewRedisRepository(log logger.Logger, cfg *config.Config, redisClient redis.UniversalClient) *redisRepository {
	return &redisRepository{log: log, cfg: cfg, redisClient: redisClient}
}

func (r *redisRepository) PutAdminToken(ctx context.Context, token string) error {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.PutAdminToken")
		defer span.End()
	}

	if token == "" {
		return errors.New("invalid input")
	}

	if err := r.redisClient.Set(ctx, r.getRedisAdminTokenKey(), token, redisAdminTokenLifetime).Err(); err != nil {
		r.log.WarnMsg("redisClient.Set", err)
		return errors.Wrap(err, "redisClient.Set")
	}

	return nil
}

func (r *redisRepository) GetAdminToken(ctx context.Context) (string, error) {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.GetAdminToken")
		defer span.End()
	}

	token := r.redisClient.Get(ctx, r.getRedisAdminTokenKey()).Val()
	if token == "" {
		return "", errors.New("redisClient.Get")
	}

	return token, nil
}

func (r *redisRepository) PutApiList(ctx context.Context, companyUuid string, apiList []*models.Api) error {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.PutApiList")
		defer span.End()
	}

	key := r.getRedisApiListKey()

	for _, api := range apiList {
		apiBytes, err := json.Marshal(api)
		if err != nil {
			r.log.WarnMsg("json.Marshal", err)
			return errors.Wrap(err, "json.Marshal")
		}

		field := r.getRedisApiListField(companyUuid, api.Uuid)

		if err := r.redisClient.HSet(ctx, key, field, apiBytes).Err(); err != nil {
			r.log.WarnMsg("redisClient.HSet", err)
			return errors.Wrap(err, "redisClient.HSet")
		}
		r.redisClient.HExpire(ctx, key, redisApiLife, field)
	}
	r.redisClient.Expire(ctx, key, redisApiLife)

	return nil
}

func (r *redisRepository) GetApiList(ctx context.Context, companyUuid string) ([]*models.Api, error) {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.GetApiList")
		defer span.End()
	}

	key := r.getRedisApiListKey()
	field := r.getRedisApiListField(companyUuid, constants.RedisAllItems)

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
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.GetApi")
		defer span.End()
	}

	if apiUuid == "" {
		return nil, errors.New("empty input")
	}

	key := r.getRedisApiListKey()
	field := r.getRedisApiListField(companyUuid, apiUuid)

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
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.DeleteApiList")
		defer span.End()
	}

	key := r.getRedisApiListKey()
	field := r.getRedisApiListField(company.Uuid, constants.RedisAllItems)

	fields, err := r.HScan(ctx, key, field)
	if err != nil {
		return err
	}

	for _, f := range fields {
		r.redisClient.HDel(ctx, key, f.Field)
	}

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
		fields, nextCursor, err := r.redisClient.HScan(ctx, key, cursor, matchPattern, 30).Result()
		if err != nil {
			return nil, errors.Wrap(err, "redisClient.HScan")
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
