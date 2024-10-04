package api

import (
	"itspay/internal/config"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type probesServer struct {
	c       *config.TechServerConfig
	pgxPool *pgxpool.Pool
	l       *zap.Logger
}

func newProbesServer(c *config.TechServerConfig, pgxPool *pgxpool.Pool, l *zap.Logger) *probesServer {
	return &probesServer{c: c, pgxPool: pgxPool, l: l}
}

func (s *probesServer) newMux() http.Handler {
	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())

	router.Methods(http.MethodGet).Path("/health").HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) { //nolint:revive
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("OK"))
		},
	)

	router.Methods(http.MethodGet).Path("/ready").HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			if err := s.pgxPool.Ping(request.Context()); err != nil {
				s.l.Error("postgres pool not ready", zap.Error(err))

				http.Error(writer, "postgres pool not ready", http.StatusInternalServerError)

				return
			}

			_, _ = writer.Write([]byte("OK"))
		},
	)

	pprofRouter := router.PathPrefix("/debug/pprof/").Subrouter()
	{
		pprofRouter.HandleFunc("/cmdline", pprof.Cmdline)
		pprofRouter.HandleFunc("/profile", pprof.Profile)
		pprofRouter.HandleFunc("/symbol", pprof.Symbol)
		pprofRouter.HandleFunc("/trace", pprof.Trace)
		pprofRouter.PathPrefix("/").HandlerFunc(pprof.Index)
	}

	return router
}

func (s *probesServer) newHTTPServer() *http.Server {
	return &http.Server{
		Addr:    s.c.Addr,
		Handler: s.newMux(),
	}
}
