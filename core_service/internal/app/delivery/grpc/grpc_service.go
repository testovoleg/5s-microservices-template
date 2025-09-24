package grpc

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	uuid "github.com/satori/go.uuid"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/commands"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/metrics"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcService struct {
	log     logger.Logger
	cfg     *config.Config
	v       *validator.Validate
	svc     *service.Service
	metrics *metrics.CoreServiceMetrics
}

func NewCoreGrpcService(log logger.Logger, cfg *config.Config, v *validator.Validate, svc *service.Service, metrics *metrics.CoreServiceMetrics) *grpcService {
	return &grpcService{log: log, cfg: cfg, v: v, svc: svc, metrics: metrics}
}

func (s *grpcService) InvoiceHandlersList(ctx context.Context, req *coreService.InvoiceHandlersListReq) (*coreService.InvoiceHandlersListRes, error) {
	s.metrics.InvoiceHandlersListGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.InvoiceHandlersList")
	defer span.End()

	command := commands.NewInvoiceHandlersListCommand()
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	res, err := s.svc.Commands.InvoiceHandlersList.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("InvoiceHandlersList.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	s.metrics.SuccessGrpcRequests.Inc()
	return &coreService.InvoiceHandlersListRes{Handlers: models.InvoiceHandlersListToGrpc(res)}, nil
}

func (s *grpcService) UpdateProduct(ctx context.Context, req *coreService.UpdateProductReq) (*coreService.UpdateProductRes, error) {
	s.metrics.UpdateProductGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.UpdateProduct")
	defer span.End()

	command := commands.NewUpdateProductCommand(req.GetProductID(), req.GetName(), req.GetDescription(), req.GetPrice(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	if err := s.svc.Commands.UpdateProduct.Handle(ctx, command); err != nil {
		s.log.WarnMsg("UpdateProduct.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	s.metrics.SuccessGrpcRequests.Inc()
	return &coreService.UpdateProductRes{ProductID: req.GetProductID()}, nil
}

func (s *grpcService) DeleteProductByID(ctx context.Context, req *coreService.DeleteProductByIdReq) (*coreService.DeleteProductByIdRes, error) {
	s.metrics.DeleteProductGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.DeleteProductByID")
	defer span.End()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	if err := s.svc.Commands.DeleteProduct.Handle(ctx, commands.NewDeleteProductCommand(productUUID)); err != nil {
		s.log.WarnMsg("DeleteProduct.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.metrics.SuccessGrpcRequests.Inc()
	return &coreService.DeleteProductByIdRes{}, nil
}

func (s *grpcService) errResponse(c codes.Code, err error) error {
	s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
