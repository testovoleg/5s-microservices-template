package middlewares

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareManager struct {
	log logger.Logger
	cfg *config.Config
}

func NewMiddlewareManager(log logger.Logger, cfg *config.Config) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg}
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
