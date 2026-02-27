package grpc

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/internal/mappers"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *grpcService) AddApi(ctx context.Context, req *coreService.AddApiReq) (*coreService.Api, error) {
	s.metrics.Get("AddApi", metrics.GRPC).Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.AddApi")
	defer span.End()

	command := mappers.NewAddApiCommandFromGrpcMessage(req)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	res, err := s.svc.Commands.AddApi.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("AddApi.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	s.metrics.Get("Success", metrics.GRPC).Inc()
	return models.ApiToGrpcMessage(res), nil
}

func (s *grpcService) GetApi(ctx context.Context, req *coreService.GetApiReq) (*coreService.GetApiRes, error) {
	s.metrics.Get("GetApi", metrics.GRPC).Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetApi")
	defer span.End()

	command := mappers.NewGetApiCommandFromGrpcMessage(req)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	res, err := s.svc.Commands.GetApi.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("GetApi.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	s.metrics.Get("Success", metrics.GRPC).Inc()
	return &coreService.GetApiRes{Api: models.ListApiToGrpcMessage(res)}, nil
}

func (s *grpcService) GetFullApi(ctx context.Context, req *coreService.GetFullApiReq) (*coreService.ApiFull, error) {
	s.metrics.Get("GetFullApi", metrics.GRPC).Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetFullApi")
	defer span.End()

	command := mappers.NewGetFullApiCommandFromGrpcMessage(req)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	res, err := s.svc.Commands.GetFullApi.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("GetFullApi.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	s.metrics.Get("Success", metrics.GRPC).Inc()
	return models.ApiFullToGrpcMessage(res), nil
}

func (s *grpcService) UpdateApi(ctx context.Context, req *coreService.UpdateApiReq) (*coreService.Api, error) {
	s.metrics.Get("UpdateApi", metrics.GRPC).Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.UpdateApi")
	defer span.End()

	command := mappers.NewUpdateApiCommandFromGrpcMessage(req)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	res, err := s.svc.Commands.UpdateApi.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("UpdateApi.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	s.metrics.Get("Success", metrics.GRPC).Inc()
	return models.ApiToGrpcMessage(res), nil
}

func (s *grpcService) DeleteApi(ctx context.Context, req *coreService.DeleteApiReq) (*emptypb.Empty, error) {
	s.metrics.Get("DeleteApi", metrics.GRPC).Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.DeleteApi")
	defer span.End()

	command := mappers.NewDeleteApiCommandFromGrpcMessage(req)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	err := s.svc.Commands.DeleteApi.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("DeleteApi.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	s.metrics.Get("Success", metrics.GRPC).Inc()
	return nil, nil
}
