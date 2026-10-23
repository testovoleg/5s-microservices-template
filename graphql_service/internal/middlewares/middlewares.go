package middlewares

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	httpUtils "github.com/testovoleg/5s-microservice-template/pkg/http_utils"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type MiddlewareManager interface {
	WriteHeadersToContext(next echo.HandlerFunc) echo.HandlerFunc
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareManager struct {
	log logger.Logger
	cfg *config.Config
}

func NewMiddlewareManager(log logger.Logger, cfg *config.Config) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg}
}

func (mw *middlewareManager) WriteHeadersToContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get(echo.HeaderAuthorization)
		if token == "" {
			token = c.Request().Header.Get(echo.HeaderWWWAuthenticate)
		}
		token = httpUtils.StripBearer(token)
		if token != "" {
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, constants.ContextKeyToken, token)
			c.SetRequest(c.Request().WithContext(ctx))
		}
		return next(c)
	}
}

func (mw *middlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return httpUtils.AccessLogMiddleware(mw.log, mw.cfg.Http.IgnoreLogUrls)(next)
}
