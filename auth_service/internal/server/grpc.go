package server

import (
	"github.com/paw1a/grpc-media-converter/auth_service/internal/repository"
	"github.com/paw1a/grpc-media-converter/auth_service/internal/service"
	"github.com/paw1a/grpc-media-converter/auth_service/pb"
	"github.com/paw1a/grpc-media-converter/auth_service/pkg/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"time"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (s *Server) newGrpcServer() (serverCloseFunc func() error, server *grpc.Server, err error) {
	l, err := net.Listen("tcp", s.cfg.GRPC.Port)
	if err != nil {
		return nil, nil, errors.Wrap(err, "net.Listen")
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
	)

	userRepo := repository.NewUserRepository(s.dbPool)
	jwt := utils.NewJwtProvider(s.cfg.JWT)

	grpcService := service.NewGrpcService(s.cfg, userRepo, jwt)
	pb.RegisterAuthServiceServer(grpcServer, grpcService)

	go func() {
		log.Printf("auth grpc server is listening on port: %s", s.cfg.GRPC.Port)
		log.Fatal(grpcServer.Serve(l))
	}()

	return l.Close, grpcServer, nil
}
