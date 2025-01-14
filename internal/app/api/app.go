package api

import (
	"context"
	"fmt"
	ratesapi "itspay/internal/api/rates"
	"itspay/internal/config"
	postgresratesdb "itspay/internal/db/rates/postgres"
	garantexrateprovider "itspay/internal/rateprovider/garantex"
	metricsrateprovider "itspay/internal/rateprovider/metrics"
	ratesservice "itspay/internal/service/rates"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	c              *config.Config
	ratesAPIServer *ratesapi.Server

	// Resources to release on shutdown below
	pgxPool *pgxpool.Pool

	// Technical things below
	l *otelzap.Logger
}

func NewApp(ctx context.Context) (*App, error) {
	zapL, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("can't initialize zap logger: %w", err)
	}

	l := otelzap.New(zapL)

	c, err := loadConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	if err := initTracerProvider(ctx, &c.OTEL); err != nil {
		return nil, fmt.Errorf("init tracer provider: %w", err)
	}

	pgxPool, err := setupPgxPool(ctx, &c.Postgres)
	if err != nil {
		return nil, fmt.Errorf("unable to setup pgx pool: %w", err)
	}

	rateDB := postgresratesdb.New(pgxPool)
	garantexRateProvider := garantexrateprovider.New()
	ratesService := ratesservice.NewSingleflightService(
		ratesservice.New(metricsrateprovider.New(garantexRateProvider), rateDB),
	)
	ratesAPIServer := ratesapi.NewServer(ratesService, l)

	return &App{
		c:              c,
		ratesAPIServer: ratesAPIServer,
		pgxPool:        pgxPool,
		l:              l,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	grpcServer := setupGRPCServer(a.ratesAPIServer, a.l)
	techHTTPServer := newTechServer(&a.c.TechServer, a.pgxPool, a.l).newHTTPServer()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", a.c.GRPC.Addr)
		if err != nil {
			return err
		}

		a.l.Info("running grpc server", zap.String("addr", lis.Addr().String()))

		return grpcServer.Serve(lis)
	})

	g.Go(func() error {
		a.l.Info("tech server start listening", zap.String("addr", a.c.TechServer.Addr))

		return techHTTPServer.ListenAndServe()
	})

	g.Go(func() error {
		<-stop

		return context.Canceled
	})

	// Goroutine below cleans up all resources
	g.Go(func() error {
		<-ctx.Done()

		a.l.Info("shutting down gracefully")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:mnd
		defer cancel()

		grpcServer.GracefulStop()

		if err := techHTTPServer.Shutdown(shutdownCtx); err != nil { //nolint:contextcheck
			a.l.Warn("failed to properly shutdown http server", zap.Error(err))
		}

		a.pgxPool.Close()

		_ = a.l.Sync()

		return nil
	})

	return g.Wait()
}
