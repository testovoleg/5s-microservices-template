package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
)

type ApiGatewayMetrics struct {
	SuccessHttpRequests             prometheus.Counter
	ErrorHttpRequests               prometheus.Counter
	InvoiceHandlersListHttpRequests prometheus.Counter
	UpdateProductHttpRequests       prometheus.Counter
	DeleteProductHttpRequests       prometheus.Counter
	GetProductByIdHttpRequests      prometheus.Counter
	SearchProductHttpRequests       prometheus.Counter
}

func NewApiGatewayMetrics(cfg *config.Config) *ApiGatewayMetrics {
	return &ApiGatewayMetrics{
		SuccessHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests_total", cfg.ServiceName),
			Help: "The total number of success http requests",
		}),
		ErrorHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests_total", cfg.ServiceName),
			Help: "The total number of error http requests",
		}),
		InvoiceHandlersListHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_invoice_handlers_list_http_requests_total", cfg.ServiceName),
			Help: "The total number of invoice handlers list http requests",
		}),
		UpdateProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of update product http requests",
		}),
		DeleteProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of delete product http requests",
		}),
		GetProductByIdHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_product_by_id_http_requests_total", cfg.ServiceName),
			Help: "The total number of get product by id http requests",
		}),
		SearchProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of search product http requests",
		}),
	}
}
