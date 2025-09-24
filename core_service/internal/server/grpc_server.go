package server

import (
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	coreGrpc "github.com/testovoleg/5s-microservice-template/core_service/internal/app/delivery/grpc"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (s *server) newCoreGrpcServer() (func() error, *grpc.Server, error) {
	l, err := net.Listen("tcp", s.cfg.GRPC.Port)
	if err != nil {
		return nil, nil, errors.Wrap(err, "net.Listen")
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			//grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
			s.im.Logger,
		),
		),
	)

	coreGrpcService := coreGrpc.NewCoreGrpcService(s.log, s.cfg, s.v, s.svc, s.metrics)
	coreService.RegisterCoreServiceServer(grpcServer, coreGrpcService)
	grpc_prometheus.Register(grpcServer)

	if s.cfg.GRPC.Development {
		reflection.Register(grpcServer)
	}

	go func() {
		s.log.Infof("Core gRPC server is listening on port: %s", s.cfg.GRPC.Port)
		// start serving grpc
		err := grpcServer.Serve(l)
		s.log.Warn(err)
		// wait 3 second, when grpc is ended to server. Is needed for log other fatals, if exists
		time.Sleep(time.Second * 3)
		// and crash error
		s.log.Fatal(err)
	}()

	return l.Close, grpcServer, nil
}
