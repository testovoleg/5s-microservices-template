package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	httpUtils "github.com/testovoleg/5s-microservice-template/pkg/http_utils"
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
	return httpUtils.AccessLogMiddleware(mw.log, mw.cfg.Http.IgnoreLogUrls)(next)
}
