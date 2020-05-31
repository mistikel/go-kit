package server

import (
	"os"
	"os/signal"
	"syscall"

	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/go-kit/kit/log"
	"github.com/mistikel/go-kit/endpoints"
	"github.com/mistikel/go-kit/pkg/imdb"
)

var IsShuttingDown = false

type Manager interface {
	Serve()
	EnableGracefulShutdown()
}

type imManager struct {
	httpServer *HTTPServer
	grpcServer *GRPCServer
}

func New() Manager {
	tracer := stdopentracing.GlobalTracer()
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	service := imdb.New(logger)
	endpoints := endpoints.NewEndpoints(service, tracer, logger)

	return &imManager{
		httpServer: NewHTTPServer(endpoints, tracer, logger),
		grpcServer: NewGRPCServer(endpoints, tracer, logger),
	}
}

func (s *imManager) Serve() {
	go s.EnableGracefulShutdown()
	go s.httpServer.Serve()
	s.grpcServer.Serve()
}

func (s *imManager) EnableGracefulShutdown() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go s.handleShutdown(signalChannel)
}

func (s *imManager) handleShutdown(ch chan os.Signal) {
	<-ch
	defer os.Exit(0)
	IsShuttingDown = true
}
