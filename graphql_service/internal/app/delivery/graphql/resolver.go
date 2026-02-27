package graph_resolvers

import (
	"github.com/go-playground/validator"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/middlewares"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
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
