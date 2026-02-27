package metrics

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type MetricsManager struct {
	Counters      []*MetricCounter
	MetricOptions MetricOptions
	log           logger.Logger
}

type MetricCounter struct {
	Name    string
	Counter prometheus.Counter
}

type MetricOptions struct {
	ServiceName      string
	MicroserviceName string
	PrometheusPath   string
	PrometheusPort   string
}

type TransportType string

const (
	HTTP    TransportType = "http"
	KAFKA   TransportType = "kafka"
	GRPC    TransportType = "grpc"
	GRAPHQL TransportType = "graphql"
)

func NewMetricsManager(log logger.Logger, serviceName, serverPath, serverPort string) *MetricsManager {
	res := &MetricsManager{
		Counters: []*MetricCounter{},
		MetricOptions: MetricOptions{
			ServiceName:      formatString(serviceName),
			MicroserviceName: formatString(constants.ShortMicroserviceName),
			PrometheusPath:   serverPath,
			PrometheusPort:   serverPort,
		},
		log: log,
	}
	return res
}

func (m *MetricsManager) NewServer(cancel context.CancelFunc, stackSize int) {
	metricsServer := echo.New()
	metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	go func() {
		metricsServer.GET(m.MetricOptions.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		m.log.Infof("Metrics server is running on port: %s", m.MetricOptions.PrometheusPort)
		if err := metricsServer.Start(m.MetricOptions.PrometheusPort); err != nil {
			m.log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()
}

func (m *MetricsManager) Get(operation string, transport TransportType) *MetricCounter {
	if operation == "" {
		return nil
	}

	name, help := m.Opts(operation, transport)
	for _, c := range m.Counters {
		if c.Name == name {
			return c
		}
	}

	counter := m.NewCounter(name, help)
	m.Counters = append(m.Counters, counter)
	prometheus.MustRegister(counter.Counter)

	return counter
}

func (m *MetricsManager) Opts(operation string, transport TransportType) (string, string) {
	operation = formatString(operation)
	if operation == "" {
		operation = "unknown"
	}

	service := fmt.Sprintf("%s_%s", m.MetricOptions.MicroserviceName, m.MetricOptions.ServiceName)

	name := fmt.Sprintf("%s_%s_%s_%s_total",
		service, operation, strings.ToLower(string(transport)), transport.Suffix())

	help := fmt.Sprintf("Total number of %s %s %s",
		strings.ReplaceAll(operation, "_", " "), strings.ToLower(string(transport)), transport.Suffix())
	return name, help
}

func (m *MetricsManager) NewCounter(n, h string) *MetricCounter {
	return &MetricCounter{
		Name: n,
		Counter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: n,
			Help: h,
		}),
	}
}

func (m *MetricCounter) Inc() {
	if m != nil {
		m.Counter.Inc()
	}
}

func (m TransportType) Suffix() string {
	switch m {
	case KAFKA:
		return "messages"
	case GRAPHQL:
		return "queries"
	}
	return "requests"
}

func (m TransportType) String() string {
	return string(m)
}

func formatString(str string) string {
	// "GetApi" → "Get_Api", "GetHTTPResponse" → "Get_HTTP_Response"
	result := regexp.MustCompile(`([a-z0-9])([A-Z])`).ReplaceAllString(str, "${1}_${2}")

	// "Get_Api" → "get_api"
	result = strings.ToLower(result)

	// "get-api" → "get_api"
	result = regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(result, "_")

	// "get__api" → "get_api"
	result = regexp.MustCompile(`_+`).ReplaceAllString(result, "_")

	// "_get_api_" → "get_api"
	result = strings.Trim(result, "_")

	return result
}
