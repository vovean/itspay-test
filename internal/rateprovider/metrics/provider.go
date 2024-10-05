package metricsrateprovider

import (
	"context"
	"itspay/internal/entity"
	"itspay/internal/rateprovider"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var getRateDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "get_rate_duration_seconds",
	Help:    "Duration of getting rate from rate provider, seconds",
	Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 3.0},
}, []string{"error"})

var _ rateprovider.RateProvider = &RateProvider{}

type RateProvider struct {
	p rateprovider.RateProvider
}

func New(p rateprovider.RateProvider) *RateProvider {
	return &RateProvider{p: p}
}

func (p *RateProvider) GetRate(ctx context.Context) (_ *entity.Rate, err error) { //nolint:nonamedreturns
	defer func(now time.Time) {
		getRateDurationSeconds.With(prometheus.Labels{
			"error": strconv.FormatBool(err != nil),
		}).Observe(time.Since(now).Seconds())
	}(time.Now())

	rate, err := p.p.GetRate(ctx)

	return rate, err
}
