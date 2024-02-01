package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
)

type CoreServiceMetrics struct {
	SuccessGrpcRequests prometheus.Counter
	ErrorGrpcRequests   prometheus.Counter

	InvoiceHandlersListGrpcRequests prometheus.Counter
	UpdateProductGrpcRequests       prometheus.Counter
	DeleteProductGrpcRequests       prometheus.Counter
	GetProductByIdGrpcRequests      prometheus.Counter
	SearchProductGrpcRequests       prometheus.Counter

	SuccessKafkaMessages prometheus.Counter
	ErrorKafkaMessages   prometheus.Counter

	CreateProductKafkaMessages prometheus.Counter
	UpdateProductKafkaMessages prometheus.Counter
	DeleteProductKafkaMessages prometheus.Counter
}

func NewCoreServiceMetrics(cfg *config.Config) *CoreServiceMetrics {
	return &CoreServiceMetrics{
		SuccessGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of success grpc requests",
		}),
		ErrorGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of error grpc requests",
		}),
		InvoiceHandlersListGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_invoice_handlers__list_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of invoice handlers list grpc requests",
		}),
		UpdateProductGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_product_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of update product grpc requests",
		}),
		DeleteProductGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_product_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of delete product grpc requests",
		}),
		GetProductByIdGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_product_by_id_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get product by id grpc requests",
		}),
		SearchProductGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_product_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of search product grpc requests",
		}),
		CreateProductKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_product_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of create product kafka messages",
		}),
		UpdateProductKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_product_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of update product kafka messages",
		}),
		DeleteProductKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_product_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of delete product kafka messages",
		}),
		SuccessKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_kafka_processed_messages_total", cfg.ServiceName),
			Help: "The total number of success kafka processed messages",
		}),
		ErrorKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_kafka_processed_messages_total", cfg.ServiceName),
			Help: "The total number of error kafka processed messages",
		}),
	}
}
