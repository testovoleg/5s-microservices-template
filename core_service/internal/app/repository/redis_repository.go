package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

const (
	redisMicroservicePrefixKey = "vk-ads-api"
	redisAdminTokenLifetime    = 3 * time.Hour
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
	if trace.SpanContextFromContext(ctx).IsValid() {
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.PutAdminToken")
		defer span.End()
	}
	if err := r.redisClient.Set(ctx, r.getRedisAdminTokenKey(), token, redisAdminTokenLifetime).Err(); err != nil {
		r.log.WarnMsg("redisClient.Set", err)
		return errors.Wrap(err, "redisClient.Set")
	}
	return nil
}

func (r *redisRepository) GetAdminToken(ctx context.Context) (string, error) {
	if trace.SpanContextFromContext(ctx).IsValid() {
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.GetAdminToken")
		defer span.End()
	}
	dataBytes, err := r.redisClient.Get(ctx, r.getRedisAdminTokenKey()).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.Get", err)
		}
		return "", errors.Wrap(err, "redisClient.Get")
	}

	return string(dataBytes), nil
}

func (r *redisRepository) PutProduct(ctx context.Context, key string, product *models.Product) {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.PutProduct")
		defer span.End()
	}

	productBytes, err := json.Marshal(product)
	if err != nil {
		r.log.WarnMsg("json.Marshal", err)
		return
	}

	if err := r.redisClient.HSetNX(ctx, r.getRedisMicroservicePrefixKey(), key, productBytes).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisMicroservicePrefixKey(), key)
}

func (r *redisRepository) GetProduct(ctx context.Context, key string) (*models.Product, error) {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.GetProduct")
		defer span.End()
	}

	productBytes, err := r.redisClient.HGet(ctx, r.getRedisMicroservicePrefixKey(), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}

	var product models.Product
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, err
	}

	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisMicroservicePrefixKey(), key)
	return &product, nil
}

func (r *redisRepository) DelProduct(ctx context.Context, key string) {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.DelProduct")
		defer span.End()
	}

	if err := r.redisClient.HDel(ctx, r.getRedisMicroservicePrefixKey(), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisMicroservicePrefixKey(), key)
}

func (r *redisRepository) DelAllProducts(ctx context.Context) {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.DelAllProducts")
		defer span.End()
	}

	if err := r.redisClient.Del(ctx, r.getRedisMicroservicePrefixKey()).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisMicroservicePrefixKey())
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
