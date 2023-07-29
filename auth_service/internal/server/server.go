package server

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paw1a/grpc-media-converter/auth_service/config"
	"github.com/paw1a/grpc-media-converter/auth_service/pkg/postgres"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	cfg    *config.Config
	dbPool *pgxpool.Pool
}

func NewServer(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	dbPool, err := postgres.NewPgxPool(s.cfg.Postgres)
	if err != nil {
		return errors.Wrap(err, "NewPgxPool")
	}
	s.dbPool = dbPool

	serverCloseFunc, grpcServer, err := s.newGrpcServer()
	if err != nil {
		return errors.Wrap(err, "newGrpcServer")
	}
	defer serverCloseFunc()

	<-ctx.Done()
	grpcServer.GracefulStop()

	return nil
}
