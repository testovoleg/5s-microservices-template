package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	httpUtils "github.com/testovoleg/5s-microservice-template/pkg/http_utils"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type authRepository struct {
	log logger.Logger
	cfg *config.Config
}

func NewAuthRepository(log logger.Logger, cfg *config.Config) *authRepository {
	return &authRepository{log: log, cfg: cfg}
}

func (r *authRepository) NewExtSession(ctx context.Context, adminToken, companyUuid, clientID string) (*models.ExtSession, error) {
	ctx, span := tracing.StartSpan(ctx, "authRepository.NewExtSession")
	defer span.End()

	requestURL := r.cfg.API.AuthApiUrl + "/ext/session"

	jsonBody, err := json.Marshal(models.ExtSessionReq{
		CompanyUuid: companyUuid,
		ClientId:    clientID,
	})
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodPost, requestURL, jsonBody, httpUtils.SendOpts{
		BearerToken: adminToken,
		Timeout:     30 * time.Second,
		ErrorPrefix: "Error was during on Auth API",
	})
	if err != nil {
		return nil, err
	}

	var body models.ExtSession
	if err = json.Unmarshal(b, &body); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	return &body, nil
}
