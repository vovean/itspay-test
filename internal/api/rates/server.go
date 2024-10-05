package ratesapi

import (
	"itspay/internal/api/rates/ratespb"
	"itspay/internal/service"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

//go:generate protoc -I ../../.. --go_out module=itspay:../../.. --go-grpc_out module=itspay:../../..  api/rates/proto/rates_service.proto

type Server struct {
	ratespb.UnimplementedRatesServiceServer
	service service.Rates
	l       *otelzap.Logger
}

func NewServer(service service.Rates, l *otelzap.Logger) *Server {
	return &Server{service: service, l: l}
}
