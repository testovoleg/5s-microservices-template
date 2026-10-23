package repository

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

func (r *redisRepository) getRedisEventListKey() string {
	return r.getRedisMicroservicePrefixKey() + ":event_list"
}

func (r *redisRepository) getRedisEventField(companyUuid, eventUuid string) string {
	return "company@" + companyUuid + ":event@" + eventUuid
}

func (r *redisRepository) PutEvent(ctx context.Context, companyUuid string, event *models.Event) error {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.PutEvent")
	defer span.End()

	if companyUuid == "" || event == nil || event.Uuid == "" {
		return httpErrors.ErrorWrap(errors.New("invalid event"), span)
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	key := r.getRedisEventListKey()
	field := r.getRedisEventField(companyUuid, event.Uuid)
	if err := r.redisClient.HSet(ctx, key, field, eventBytes).Err(); err != nil {
		return httpErrors.ErrorWrap(err, span, "redisClient.HSet")
	}
	r.redisClient.HExpire(ctx, key, redisEventLife, field)

	return nil
}

func (r *redisRepository) GetEvent(ctx context.Context, companyUuid, eventUuid string) (*models.Event, error) {
	ctx, span := tracing.StartSpan(ctx, "redisRepository.GetEvent")
	defer span.End()

	eventBytes, err := r.redisClient.HGet(ctx, r.getRedisEventListKey(), r.getRedisEventField(companyUuid, eventUuid)).Bytes()
	if err != nil {
		return &models.Event{Uuid: eventUuid, Status: models.EVENT_STATUS_UNKNOWN}, nil
	}

	var event models.Event
	if err := json.Unmarshal(eventBytes, &event); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	return &event, nil
}
