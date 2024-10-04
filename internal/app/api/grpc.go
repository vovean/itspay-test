package api

import (
	"context"
	"fmt"
	ratesapi "itspay/internal/api/rates"
	"itspay/internal/api/rates/ratespb"
	grpcmetrics "itspay/internal/utils/grpckit/metrics"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// interceptorLogger adapts zap logger to interceptor logger.
func interceptorLogger(l *otelzap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2) //nolint:mnd

		for i := 0; i < len(fields); i += 2 {
			key, ok := fields[i].(string)
			if !ok {
				continue
			}

			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key, v))
			case int:
				f = append(f, zap.Int(key, v))
			case bool:
				f = append(f, zap.Bool(key, v))
			default:
				f = append(f, zap.Any(key, v))
			}
		}

		logger := l.Ctx(ctx).WithOptions(zap.AddCallerSkip(1))

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg, f...)
		case logging.LevelInfo:
			logger.Info(msg, f...)
		case logging.LevelWarn:
			logger.Warn(msg, f...)
		case logging.LevelError:
			logger.Error(msg, f...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func setupGRPCServer(server *ratesapi.Server, l *otelzap.Logger) *grpc.Server {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(l), loggingOpts...),
			grpcmetrics.ServerMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(interceptorLogger(l), loggingOpts...),
			grpcmetrics.ServerMetrics.StreamServerInterceptor(),
		),
	)

	ratespb.RegisterRatesServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	return grpcServer
}
