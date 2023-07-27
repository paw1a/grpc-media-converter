package server

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/paw1a/grpc-media-converter/storage_service/config"
	minio2 "github.com/paw1a/grpc-media-converter/storage_service/pkg/minio"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	cfg         *config.Config
	minioClient *minio.Client
}

func NewServer(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	minioClient, err := minio2.NewMinioClient(s.cfg.Minio)
	if err != nil {
		return errors.Wrap(err, "NewMinioClient")
	}
	log.Printf("minio connected: %v\n", minioClient.EndpointURL())
	s.minioClient = minioClient

	serverCloseFunc, grpcServer, err := s.newGrpcServer()
	if err != nil {
		return errors.Wrap(err, "newGrpcServer")
	}
	defer serverCloseFunc()

	<-ctx.Done()
	grpcServer.GracefulStop()

	return nil
}
