package grpc

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/internal/mappers"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"google.golang.org/grpc/codes"
)

func (s *grpcService) EventGet(ctx context.Context, req *coreService.EventGetReq) (*coreService.EventGetRes, error) {
	s.metrics.Get("EventGet", metrics.GRPC).Inc()

	ctx, span := tracing.StartSpan(ctx, "grpcService.EventGet")
	defer span.End()

	command := mappers.NewEventGetCommandFromGrpcMessage(req)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	res, err := s.svc.Commands.CommonMethods.EventGet(ctx, command)
	if err != nil {
		s.log.WarnMsg("CommonMethods.EventGet", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	s.metrics.Get("Success", metrics.GRPC).Inc()
	return &coreService.EventGetRes{
		Data: models.EventToGrpcMessage(res),
	}, nil
}
