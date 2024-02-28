package client

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	"github.com/testovoleg/5s-microservice-template/pkg/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	backoffLinear  = 100 * time.Millisecond
	backoffRetries = 3
)

func NewCoreServiceConn(ctx context.Context, cfg *config.Config, im interceptors.InterceptorManager) (*grpc.ClientConn, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backoffLinear)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
		grpc_retry.WithMax(backoffRetries),
	}

	coreServiceConn, err := grpc.DialContext(
		ctx,
		cfg.Grpc.CoreServicePort,
		grpc.WithUnaryInterceptor(im.ClientRequestLoggerInterceptor()),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "grpc.DialContext")
	}

	return coreServiceConn, nil
}
