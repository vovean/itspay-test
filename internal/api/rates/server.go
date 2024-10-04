package ratesapi

import (
	"context"
	"itspay/internal/api/rates/ratespb"
	"itspay/internal/service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate protoc -I ../../.. --go_out module=itspay:../../.. --go-grpc_out module=itspay:../../..  api/rates/proto/rates_service.proto

type Server struct {
	ratespb.UnimplementedRatesServiceServer
	service service.Rates
}

func NewServer(service service.Rates) *Server {
	return &Server{service: service}
}

func (s *Server) GetRate(ctx context.Context, _ *ratespb.GetRateRequest) (*ratespb.GetRateResponse, error) {
	rate, err := s.service.GetRate(ctx)
	if err != nil {
		return nil, err
	}

	return &ratespb.GetRateResponse{
		Ask:        rate.Ask.String(),
		Bid:        rate.Bid.String(),
		ReceivedAt: timestamppb.New(rate.ReceivedAt),
	}, nil
}
