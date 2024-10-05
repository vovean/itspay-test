package ratesapi

import (
	"context"
	"itspay/internal/api/rates/ratespb"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errInternal = status.Error(codes.Internal, "try again later")

func (s *Server) GetRate(ctx context.Context, _ *ratespb.GetRateRequest) (*ratespb.GetRateResponse, error) {
	rate, err := s.service.GetRate(ctx)
	if err != nil {
		s.l.Ctx(ctx).Error("error getting rate", zap.Error(err))

		return nil, errInternal
	}

	return &ratespb.GetRateResponse{
		Ask:        rate.Ask.String(),
		Bid:        rate.Bid.String(),
		ReceivedAt: timestamppb.New(rate.ReceivedAt),
	}, nil
}
