package httpUtils

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

// PathContainsSkipper skips Echo middleware when the request path contains any of parts.
func PathContainsSkipper(parts ...string) middleware.Skipper {
	return func(c echo.Context) bool {
		return pathContains(c.Request().URL.Path, parts...)
	}
}

// ProbeSkipper skips tracing/logging on k8s probes and Prometheus scrapes.
func ProbeSkipper(c echo.Context) bool {
	return pathContains(c.Request().URL.Path, "/health", "/ready", "/live", "metrics")
}

func pathContains(path string, parts ...string) bool {
	for _, p := range parts {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

// URIIgnored reports whether requestURI matches any configured ignore fragment.
func URIIgnored(requestURI string, uriList []string) bool {
	return pathContains(requestURI, uriList...)
}

// ResolveHTTPTraceID returns the request trace id (echo context or span) and sets X-Trace-ID.
func ResolveHTTPTraceID(c echo.Context) string {
	if traceID, _ := c.Get("traceid").(string); traceID != "" {
		if c.Response().Header().Get("X-Trace-ID") == "" {
			c.Response().Header().Set("X-Trace-ID", traceID)
		}
		return traceID
	}
	if traceID := tracing.TraceIDFromContext(c.Request().Context()); traceID != "" {
		c.Set("traceid", traceID)
		c.Response().Header().Set("X-Trace-ID", traceID)
		return traceID
	}
	return ""
}

// AccessLogMiddleware writes HTTP access logs after the handler, skipping ignoreURIs.
func AccessLogMiddleware(log logger.Logger, ignoreURIs []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			if !URIIgnored(c.Request().RequestURI, ignoreURIs) {
				req := c.Request()
				res := c.Response()
				log.HttpMiddlewareAccessLogger(req.Method, req.URL.String(), res.Status, res.Size, time.Since(start), ResolveHTTPTraceID(c))
			}
			return err
		}
	}
}
