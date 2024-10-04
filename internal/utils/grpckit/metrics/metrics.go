package grpcmetrics

import (
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var durationSecondsHistogramBuckets = []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 3.0}

var (
	ClientMetrics = grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets(durationSecondsHistogramBuckets),
		),
	)
	ServerMetrics = grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets(durationSecondsHistogramBuckets),
		),
	)
)

func init() {
	prometheus.MustRegister(ClientMetrics)
	prometheus.MustRegister(ServerMetrics)
}
