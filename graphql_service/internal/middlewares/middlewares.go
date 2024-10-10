package middlewares

import (
	"context"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
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
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}
		if token != "" {
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, constants.ContextKeyToken, token)
			c.SetRequest(c.Request().WithContext(ctx))
		}
		return next(c)
	}
}

func (mw *middlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start)
		traceID, _ := ctx.Get("traceid").(string)

		if !mw.checkIgnoredURI(ctx.Request().RequestURI, mw.cfg.Http.IgnoreLogUrls) {
			mw.log.HttpMiddlewareAccessLogger(req.Method, req.URL.String(), status, size, s, traceID)
		}

		return err
	}
}

func (mw *middlewareManager) checkIgnoredURI(requestURI string, uriList []string) bool {
	for _, s := range uriList {
		if strings.Contains(requestURI, s) {
			return true
		}
	}
	return false
}
