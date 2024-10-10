package repository

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

const (
	redisProductPrefixKey = "reader:product"
)

type redisRepository struct {
	log         logger.Logger
	cfg         *config.Config
	redisClient redis.UniversalClient
}

func NewRedisRepository(log logger.Logger, cfg *config.Config, redisClient redis.UniversalClient) *redisRepository {
	return &redisRepository{log: log, cfg: cfg, redisClient: redisClient}
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

	if err := r.redisClient.HSetNX(ctx, r.getRedisProductPrefixKey(), key, productBytes).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisProductPrefixKey(), key)
}

func (r *redisRepository) GetProduct(ctx context.Context, key string) (*models.Product, error) {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.GetProduct")
		defer span.End()
	}

	productBytes, err := r.redisClient.HGet(ctx, r.getRedisProductPrefixKey(), key).Bytes()
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

	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisProductPrefixKey(), key)
	return &product, nil
}

func (r *redisRepository) DelProduct(ctx context.Context, key string) {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.DelProduct")
		defer span.End()
	}

	if err := r.redisClient.HDel(ctx, r.getRedisProductPrefixKey(), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisProductPrefixKey(), key)
}

func (r *redisRepository) DelAllProducts(ctx context.Context) {
	if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "redisRepository.DelAllProducts")
		defer span.End()
	}

	if err := r.redisClient.Del(ctx, r.getRedisProductPrefixKey()).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisProductPrefixKey())
}

func (r *redisRepository) getRedisProductPrefixKey() string {
	if r.cfg.ServiceSettings.RedisProductPrefixKey != "" {
		return r.cfg.ServiceSettings.RedisProductPrefixKey
	}

	return redisProductPrefixKey
}
