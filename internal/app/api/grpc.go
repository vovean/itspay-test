package api

import (
	ratesapi "itspay/internal/api/rates"
	"itspay/internal/api/rates/ratespb"
	grpcmetrics "itspay/internal/utils/grpckit/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func setupGRPCServer(server *ratesapi.Server) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmetrics.ServerMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcmetrics.ServerMetrics.StreamServerInterceptor(),
		),
	)

	ratespb.RegisterRatesServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	return grpcServer
}
