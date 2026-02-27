package grpc

import (
	"github.com/go-playground/validator"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/service"

	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcService struct {
	log     logger.Logger
	cfg     *config.Config
	v       *validator.Validate
	svc     *service.CoreService
	metrics *metrics.MetricsManager
}

func NewCoreGrpcService(log logger.Logger, cfg *config.Config, v *validator.Validate, svc *service.CoreService, metrics *metrics.MetricsManager) *grpcService {
	return &grpcService{log: log, cfg: cfg, v: v, svc: svc, metrics: metrics}
}

func (s *grpcService) errResponse(c codes.Code, err error) error {
	s.metrics.Get("Error", metrics.GRPC).Inc()

	return status.Error(c, err.Error())
}
