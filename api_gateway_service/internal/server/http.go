package server

import (
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/testovoleg/5s-microservice-template/docs"
	httpUtils "github.com/testovoleg/5s-microservice-template/pkg/http_utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

const (
	maxHeaderBytes   = 1 << 20
	stackSize        = 1 << 10 // 1 KB
	defaultBodyLimit = "2M"
	readTimeout      = 15 * time.Second
	writeTimeout     = 15 * time.Second
	gzipLevel        = 5
)

func (s *server) runHttpServer() error {
	s.mapRoutes()

	s.echo.Server.ReadTimeout = readTimeout
	s.echo.Server.WriteTimeout = writeTimeout
	s.echo.Server.MaxHeaderBytes = maxHeaderBytes

	return s.echo.Start(s.cfg.Http.Port)
}

func (s *server) mapRoutes() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = s.cfg.Http.Title
	docs.SwaggerInfo.Description = "API Gateway microservices."
	docs.SwaggerInfo.BasePath = s.cfg.Http.BasePath

	// otelecho must be outermost so access logs see TRACEID from request context.
	s.echo.Use(otelecho.Middleware(s.cfg.OTL.ServiceName, otelecho.WithSkipper(httpUtils.ProbeSkipper)))
	s.echo.Use(s.mw.RequestLoggerMiddleware)
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:   gzipLevel,
		Skipper: httpUtils.PathContainsSkipper("swagger"),
	}))
	bodyLimit := s.cfg.Http.BodyLimit
	if bodyLimit == "" {
		bodyLimit = defaultBodyLimit
	}
	s.echo.Use(middleware.BodyLimit(bodyLimit))
}
