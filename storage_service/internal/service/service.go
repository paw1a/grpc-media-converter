package service

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/paw1a/grpc-media-converter/storage_service/config"
	"github.com/paw1a/grpc-media-converter/storage_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcService struct {
	pb.UnimplementedStorageServiceServer
	cfg         *config.Config
	minioClient *minio.Client
}

func NewGrpcService(cfg *config.Config, minioClient *minio.Client) *GrpcService {
	return &GrpcService{cfg: cfg, minioClient: minioClient}
}

func (s *GrpcService) UploadFile(stream pb.StorageService_UploadFileServer) error {
	return nil
}

func (s *GrpcService) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	return &pb.DownloadFileResponse{}, nil
}

func (s *GrpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
