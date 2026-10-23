package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	httpUtils "github.com/testovoleg/5s-microservice-template/pkg/http_utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

const (
	maxHeaderBytes = 1 << 20
	stackSize      = 1 << 10 // 1 KB
	bodyLimit      = "2M"
	readTimeout    = 15 * time.Second
	writeTimeout   = 15 * time.Second
	gzipLevel      = 5
)

func (s *server) runGraphQLServer() error {
	s.mapRoutes()

	s.echo.Server.ReadTimeout = readTimeout
	s.echo.Server.WriteTimeout = writeTimeout
	s.echo.Server.MaxHeaderBytes = maxHeaderBytes

	return s.echo.Start(s.cfg.Http.Port)
}

func (s *server) mapRoutes() {
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
	s.echo.Use(middleware.BodyLimit(bodyLimit))

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5001",
			"https://api.5systems.ru",
			"https://studio.apollographql.com",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
			"Apollo-Require-Preflight",
			"X-Trace-ID",
		},
		AllowCredentials: true,
	}))

	s.echo.Use(s.mw.WriteHeadersToContext)
	// s.echo.Use(s.mw.CheckAuthorization)

	g := s.echo.Group(s.cfg.Http.GraphQLPath)

	playgroundHandler := playground.Handler("GraphQL", "/query")

	g.POST("/query", func(c echo.Context) error {
		s.gql.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	g.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	g.GET("/sandbox", func(c echo.Context) error {
		c.HTML(200, fmt.Sprintf(sandboxHTML, s.cfg.Resources.GRAPHQL_QUERY))
		return nil
	})
}
