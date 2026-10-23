package interceptors

import (
	"context"
	"time"

	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptorManager interface {
	Logger(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error)
	ClientRequestLoggerInterceptor() func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error
}

// InterceptorManager struct
type interceptorManager struct {
	logger logger.Logger
}

// NewInterceptorManager InterceptorManager constructor
func NewInterceptorManager(logger logger.Logger) *interceptorManager {
	return &interceptorManager{logger: logger}
}

// Logger Interceptor
func (im *interceptorManager) Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)

	// Start a server span here so TRACEID is always present in access logs
	// (handlers may start nested spans, but this ctx keeps the trace id).
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, info.FullMethod)
	defer span.End()

	reply, err := handler(ctx, req)
	im.logger.GrpcMiddlewareAccessLogger(
		info.FullMethod,
		time.Since(start),
		md,
		err,
		tracing.TraceIDFromContext(ctx),
	)
	return reply, err
}

// ClientRequestLoggerInterceptor gRPC client interceptor
func (im *interceptorManager) ClientRequestLoggerInterceptor() func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		md, _ := metadata.FromIncomingContext(ctx)
		im.logger.GrpcClientInterceptorLogger(
			method,
			req,
			reply,
			time.Since(start),
			md,
			err,
			tracing.TraceIDFromContext(ctx),
		)
		return err
	}
}

// TracingClientInterceptor injects the caller's current span context into the
// outgoing gRPC metadata on every unary call, so the callee can extract it and
// continue the same trace (see server-side counterpart: interceptorManager.Logger,
// which calls tracing.StartGrpcServerTracerSpan for every incoming RPC).
//
// This must be registered once per gRPC client connection instead of being
// called manually from each command handler - a manual per-call approach is
// easy to forget on new methods and silently breaks trace continuity.
func TracingClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, trace.SpanContextFromContext(ctx))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

const DefaultGRPCClientTimeout = 30 * time.Second

// ClientTimeoutInterceptor sets a per-RPC deadline when the call context has none.
func ClientTimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	if timeout <= 0 {
		timeout = DefaultGRPCClientTimeout
	}
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
