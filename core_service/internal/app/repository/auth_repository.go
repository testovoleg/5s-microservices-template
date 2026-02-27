package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
	"go.opentelemetry.io/otel/trace"
)

type authRepository struct {
	log logger.Logger
	cfg *config.Config
}

func NewAuthRepository(log logger.Logger, cfg *config.Config) *authRepository {
	return &authRepository{log: log, cfg: cfg}
}

func (r *authRepository) NewExtSession(ctx context.Context, adminToken, companyUuid, clientID string) (*models.ExtSession, error) {
	var span trace.Span
	ctx, span = tracing.StartSpan(ctx, "authRepository.NewExtSession")
	defer span.End()
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestURL := r.cfg.API.AuthApiUrl + "/ext/session"

	reqBody := map[string]interface{}{}
	reqBody["company_uuid"] = companyUuid
	reqBody["client_id"] = clientID
	// reqBody["ttl"] = r.cfg.ServiceSettings.ExtTokenTelegramTTL

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
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
		return nil, errors.Wrap(err, "io.ReadAll")
	}

	var body models.ExtSession

	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	if res.StatusCode != http.StatusOK {
		return nil, utils.NewError(res, span, "Error was during on Auth API")
	}

	return &body, nil
}
