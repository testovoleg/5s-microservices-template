package graph_resolvers

import (
	"context"
	"errors"

	"github.com/go-playground/validator"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/app/service"
	model "github.com/testovoleg/5s-microservice-template/graphql_service/internal/graph_model"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/middlewares"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	log     logger.Logger
	mw      middlewares.MiddlewareManager
	cfg     *config.Config
	bs      *service.GraphQLService
	v       *validator.Validate
	metrics *metrics.MetricsManager
}

func NewResolverHandlers(
	log logger.Logger,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	bs *service.GraphQLService,
	v *validator.Validate,
	metrics *metrics.MetricsManager,
) *Resolver {
	return &Resolver{log: log, mw: mw, cfg: cfg, bs: bs, v: v, metrics: metrics}
}

// resolveAccessToken fills params.AccessToken from the Authorization header
// (via WriteHeadersToContext). Client-supplied accessToken in the GraphQL body is ignored.
func resolveAccessToken(ctx context.Context, params *model.GeneralParamsInput) error {
	if params == nil {
		return errors.New("params is required")
	}
	token, ok := ctx.Value(constants.ContextKeyToken).(string)
	if !ok || token == "" {
		return errors.New("can't find token")
	}
	params.AccessToken = &token
	return nil
}

// handleGraphQLError records a failed operation on the span, logs it and
// bumps the error metric. It returns err unchanged so callers can inline it.
func (r *Resolver) handleGraphQLError(span trace.Span, operation string, err error) error {
	r.log.WarnMsg(operation, err)
	r.metrics.Get("Error", metrics.GRAPHQL).Inc()
	return tracing.TraceErr(span, "operation in gql failed", err)
}
